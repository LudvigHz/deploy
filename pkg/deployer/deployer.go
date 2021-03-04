package deployer

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aymerick/raymond"
	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/jsonpb"
	"github.com/navikt/deployment/pkg/grpc/interceptor/apikey"
	"github.com/navikt/deployment/pkg/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	yamlv2 "gopkg.in/yaml.v2"
)

type TemplateVariables map[string]interface{}

type ActionsFormatter struct{}

type ExitCode int

const (
	DefaultRef           = "master"
	DefaultOwner         = "navikt"
	DefaultDeployServer  = "deploy-grpc.nais.io:9090"
	DefaultDeployTimeout = time.Minute * 10

	ResourceRequiredMsg = "at least one Kubernetes resource is required to make sense of the deployment"
	APIKeyRequiredMsg   = "API key required"
	MalformedURLMsg     = "wrong format of deployment server URL"
	ClusterRequiredMsg  = "cluster required; see https://doc.nais.io/clusters"
	MalformedAPIKeyMsg  = "API key must be a hex encoded string"
)

// Kept separate to avoid skewing exit codes
const (
	ExitSuccess ExitCode = iota
	ExitDeploymentFailure
	ExitDeploymentError
	ExitDeploymentInactive
	ExitNoDeployment
	ExitUnavailable
	ExitInvocationFailure
	ExitInternalError
	ExitTemplateError
	ExitTimeout
)

type Deployer struct {
	Client       *http.Client
	DeployServer string
}

func (d *Deployer) Run(cfg Config) (ExitCode, error) {
	setupLogging(cfg.Actions, cfg.Quiet)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if err := validate(cfg); err != nil {
		if !cfg.DryRun {
			return ExitInvocationFailure, err
		}

		log.Warnf("Config did not pass validation: %s", err)
	}

	var err error
	var templateVariables = make(TemplateVariables)

	if len(cfg.VariablesFile) > 0 {
		templateVariables, err = templateVariablesFromFile(cfg.VariablesFile)
		if err != nil {
			return ExitInvocationFailure, fmt.Errorf("load template variables: %s", err)
		}
	}

	if len(cfg.Variables) > 0 {
		templateOverrides := templateVariablesFromSlice(cfg.Variables)
		for key, val := range templateOverrides {
			if oldval, ok := templateVariables[key]; ok {
				log.Warnf("Overwriting template variable '%s'; previous value was '%v'", key, oldval)
			}
			log.Infof("Setting template variable '%s' to '%v'", key, val)
			templateVariables[key] = val
		}
	}

	resources := make([]json.RawMessage, 0)

	for i, path := range cfg.Resource {
		parsed, err := MultiDocumentFileAsJSON(path, templateVariables)
		if err != nil {
			if cfg.PrintPayload {
				errStr := err.Error()[len(path)+2:]
				line, er := detectErrorLine(errStr)
				if er == nil {
					ctx := errorContext(string(resources[i]), line, 7)
					for _, l := range ctx {
						fmt.Println(l)
					}
				}
			}
			return ExitTemplateError, err
		}
		resources = append(resources, parsed...)
	}

	if len(cfg.Team) == 0 {
		log.Infof("Team not explicitly specified; attempting auto-detection...")
		for i, path := range cfg.Resource {
			team := detectTeam(resources[i])
			if len(team) > 0 {
				log.Infof("Detected team '%s' in path %s", team, path)
				cfg.Team = team
				break
			}
		}

		if len(cfg.Team) == 0 {
			return ExitInvocationFailure, fmt.Errorf("no team specified, and unable to auto-detect from nais.yaml")
		}
	}

	if len(cfg.Environment) == 0 {
		log.Infof("Environment not explicitly specified; attempting auto-detection...")

		namespaces := make(map[string]interface{})
		cfg.Environment = cfg.Cluster

		for i := range cfg.Resource {
			namespace := detectNamespace(resources[i])
			namespaces[namespace] = new(interface{})
		}

		if len(namespaces) == 1 {
			for namespace := range namespaces {
				if len(namespace) != 0 {
					cfg.Environment = fmt.Sprintf("%s:%s", cfg.Cluster, namespace)
				}
			}
		}

		log.Infof("Detected environment '%s'", cfg.Environment)
	}

	allResources, err := wrapResources(resources)
	if err != nil {
		return ExitInvocationFailure, err
	}

	kube, err := pb.KubernetesFromJSONResources(allResources)
	if err != nil {
		return ExitInvocationFailure, err
	}

	deadline, _ := ctx.Deadline()
	deployRequest := &pb.DeploymentRequest{
		Cluster:           cfg.Cluster,
		Deadline:          pb.TimeAsTimestamp(deadline),
		GitRefSha:         cfg.Ref,
		GithubEnvironment: cfg.Environment,
		Kubernetes:        kube,
		Repository: &pb.GithubRepository{
			Owner: cfg.Owner,
			Name:  cfg.Repository,
		},
		Team: cfg.Team,
		Time: pb.TimeAsTimestamp(time.Now()),
	}

	if cfg.PrintPayload {
		marsh := jsonpb.Marshaler{Indent: "  "}
		err = marsh.Marshal(os.Stdout, deployRequest)
		if err != nil {
			log.Errorf("print payload: %s", err)
		}
	}

	if cfg.DryRun {
		return ExitSuccess, nil
	}

	dialOptions := make([]grpc.DialOption, 0)
	dialOptions = append(dialOptions, grpc.WithBlock())

	if !cfg.GrpcUseTLS {
		dialOptions = append(dialOptions, grpc.WithInsecure())
	} else {
		tlsOpts := &tls.Config{}
		cred := credentials.NewTLS(tlsOpts)
		if err != nil {
			return ExitInvocationFailure, fmt.Errorf("gRPC configured to use TLS, but system-wide CA certificate bundle cannot be loaded")
		}
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(cred))
	}

	if cfg.GrpcAuthentication {
		decoded, err := hex.DecodeString(cfg.APIKey)
		if err != nil {
			return ExitInvocationFailure, fmt.Errorf("%s: %s", MalformedAPIKeyMsg, err)
		}
		intercept := &apikey_interceptor.ClientInterceptor{
			APIKey:     decoded,
			RequireTLS: cfg.GrpcUseTLS,
			Team:       cfg.Team,
		}
		dialOptions = append(dialOptions, grpc.WithPerRPCCredentials(intercept))
	}

	log.Infof("Connecting to NAIS deploy at %s...", cfg.DeployServerURL)
	grpcConnection, err := grpc.DialContext(ctx, cfg.DeployServerURL, dialOptions...)
	if err != nil {
		return ExitUnavailable, fmt.Errorf("connecting to NAIS deploy: %s", err)
	}
	log.Infof("Connected to NAIS deploy; sending deployment request...")
	defer grpcConnection.Close()

	grpcClient := pb.NewDeployClient(grpcConnection)
	status, err := grpcClient.Deploy(ctx, deployRequest)
	if err != nil {
		return ExitUnavailable, err
	}

	log.Infof("Deployment request sent.")
	log.Infof("deployment: %s: %s", status.GetState(), status.GetMessage())

	if status.GetState().Finished() {
		return exitStatus(status), nil
	}

	if !cfg.Wait {
		return ExitSuccess, nil
	}

	deployRequest.ID = status.GetRequest().GetID()
	stream, err := grpcClient.Status(ctx, deployRequest)
	if err != nil {
		return ExitUnavailable, err
	}

	for ctx.Err() == nil {
		status, err = stream.Recv()
		if err != nil {
			return ExitUnavailable, err
		}
		log.Infof("deployment: %s: %s", status.GetState(), status.GetMessage())
		if status.GetState().Finished() {
			return exitStatus(status), nil
		}
	}

	return ExitTimeout, nil
}

func exitStatus(status *pb.DeploymentStatus) ExitCode {
	switch status.GetState() {
	default:
		return ExitSuccess
	case pb.DeploymentState_error:
		return ExitDeploymentError
	case pb.DeploymentState_failure:
		return ExitDeploymentFailure
	case pb.DeploymentState_inactive:
		return ExitDeploymentInactive
	}
}

func setupLogging(actions, quiet bool) {
	log.SetOutput(os.Stderr)

	if actions {
		log.SetFormatter(&ActionsFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp:          true,
			TimestampFormat:        time.RFC3339Nano,
			DisableLevelTruncation: true,
		})
	}

	if quiet {
		log.SetLevel(log.ErrorLevel)
	}
}

func detectTeam(resource json.RawMessage) string {
	type teamMeta struct {
		Metadata struct {
			Labels struct {
				Team string `json:"team"`
			} `json:"labels"`
		} `json:"metadata"`
	}
	buf := &teamMeta{}
	err := json.Unmarshal(resource, buf)

	if err != nil {
		return ""
	}

	return buf.Metadata.Labels.Team
}

func detectNamespace(resource json.RawMessage) string {
	type namespaceMeta struct {
		Metadata struct {
			Namespace string `json:"namespace"`
		} `json:"metadata"`
	}
	buf := &namespaceMeta{}
	err := json.Unmarshal(resource, buf)

	if err != nil {
		return ""
	}

	return buf.Metadata.Namespace
}

// Wrap JSON resources in a JSON array.
func wrapResources(resources []json.RawMessage) (json.RawMessage, error) {
	return json.Marshal(resources)
}

func templatedFile(data []byte, ctx TemplateVariables) ([]byte, error) {
	if len(ctx) == 0 {
		return data, nil
	}
	template, err := raymond.Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse template file: %s", err)
	}

	output, err := template.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("execute template: %s", err)
	}

	return []byte(output), nil
}

func templateVariablesFromFile(path string) (TemplateVariables, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s: open file: %s", path, err)
	}

	vars := TemplateVariables{}
	err = yaml.Unmarshal(file, &vars)

	return vars, err
}

func templateVariablesFromSlice(vars []string) TemplateVariables {
	tv := TemplateVariables{}
	for _, keyval := range vars {
		tokens := strings.SplitN(keyval, "=", 2)
		switch len(tokens) {
		case 2: // KEY=VAL
			tv[tokens[0]] = tokens[1]
		case 1: // KEY
			tv[tokens[0]] = true
		default:
			continue
		}
	}

	return tv
}

func detectErrorLine(e string) (int, error) {
	var line int
	_, err := fmt.Sscanf(e, "yaml: line %d:", &line)
	return line, err
}

func errorContext(content string, line int, around int) []string {
	ctx := make([]string, 0)
	lines := strings.Split(content, "\n")
	format := "%03d: %s"
	for l := range lines {
		ctx = append(ctx, fmt.Sprintf(format, l+1, lines[l]))
		if l+1 == line {
			helper := "     " + strings.Repeat("^", len(lines[l])) + " <--- error near this line"
			ctx = append(ctx, helper)
		}
	}
	return ctx
}

func MultiDocumentFileAsJSON(path string, ctx TemplateVariables) ([]json.RawMessage, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s: open file: %s", path, err)
	}

	templated, err := templatedFile(file, ctx)
	if err != nil {
		errMsg := strings.ReplaceAll(err.Error(), "\n", ": ")
		return nil, fmt.Errorf("%s: %s", path, errMsg)
	}

	var content interface{}
	messages := make([]json.RawMessage, 0)

	decoder := yamlv2.NewDecoder(bytes.NewReader(templated))
	for {
		err = decoder.Decode(&content)
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return nil, err
		}

		rawdocument, err := yamlv2.Marshal(content)
		if err != nil {
			return nil, err
		}

		data, err := yaml.YAMLToJSON(rawdocument)
		if err != nil {
			errMsg := strings.ReplaceAll(err.Error(), "\n", ": ")
			return nil, fmt.Errorf("%s: %s", path, errMsg)
		}

		messages = append(messages, data)
	}

	return messages, err
}

func (a *ActionsFormatter) Format(e *log.Entry) ([]byte, error) {
	buf := &bytes.Buffer{}
	switch e.Level {
	case log.ErrorLevel:
		buf.WriteString("::error::")
	case log.WarnLevel:
		buf.WriteString("::warn::")
	default:
		buf.WriteString("[")
		buf.WriteString(e.Time.Format(time.RFC3339Nano))
		buf.WriteString("] ")
	}
	buf.WriteString(e.Message)
	buf.WriteRune('\n')
	return buf.Bytes(), nil
}

func validate(cfg Config) error {
	if len(cfg.Resource) == 0 {
		return fmt.Errorf(ResourceRequiredMsg)
	}

	_, err := url.Parse(cfg.DeployServerURL)
	if err != nil {
		return fmt.Errorf("%s: %s", MalformedURLMsg, err)
	}

	if len(cfg.Cluster) == 0 {
		return fmt.Errorf(ClusterRequiredMsg)
	}

	if len(cfg.APIKey) == 0 {
		return fmt.Errorf(APIKeyRequiredMsg)
	}

	return nil
}
