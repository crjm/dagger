package dagger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vito/progrock"

	"dagger.io/dagger/querybuilder"
)

// assertNotNil panic if the given value is nil.
// This function is used to validate that input with pointer type are not nil.
// See https://github.com/dagger/dagger/issues/5696 for more context.
func assertNotNil(argName string, value any) {
	// We use reflect because just comparing value to nil is not working since
	// the value is wrapped into a type when passed as parameter.
	// E.g., nil become (*dagger.File)(nil).
	if reflect.ValueOf(value).IsNil() {
		panic(fmt.Sprintf("unexpected nil pointer for argument %q", argName))
	}
}

type DaggerObject querybuilder.GraphQLMarshaller

// getCustomError parses a GraphQL error into a more specific error type.
func getCustomError(err error) error {
	var gqlErr *gqlerror.Error

	if !errors.As(err, &gqlErr) {
		return nil
	}

	ext := gqlErr.Extensions

	typ, ok := ext["_type"].(string)
	if !ok {
		return nil
	}

	if typ == "EXEC_ERROR" {
		e := &ExecError{
			original: err,
		}
		if code, ok := ext["exitCode"].(float64); ok {
			e.ExitCode = int(code)
		}
		if args, ok := ext["cmd"].([]interface{}); ok {
			cmd := make([]string, len(args))
			for i, v := range args {
				cmd[i] = v.(string)
			}
			e.Cmd = cmd
		}
		if stdout, ok := ext["stdout"].(string); ok {
			e.Stdout = stdout
		}
		if stderr, ok := ext["stderr"].(string); ok {
			e.Stderr = stderr
		}
		return e
	}

	return nil
}

// ExecError is an API error from an exec operation.
type ExecError struct {
	original error
	Cmd      []string
	ExitCode int
	Stdout   string
	Stderr   string
}

func (e *ExecError) Error() string {
	// As a default when just printing the error, include the stdout
	// and stderr for visibility
	msg := e.Message()
	if strings.TrimSpace(e.Stdout) != "" {
		msg += "\nStdout:\n" + e.Stdout
	}
	if strings.TrimSpace(e.Stderr) != "" {
		msg += "\nStderr:\n" + e.Stderr
	}
	return msg
}

func (e *ExecError) Message() string {
	return e.original.Error()
}

func (e *ExecError) Unwrap() error {
	return e.original
}

// The `CacheVolumeID` scalar type represents an identifier for an object of type CacheVolume.
type CacheVolumeID string

// The `ContainerID` scalar type represents an identifier for an object of type Container.
type ContainerID string

// The `CurrentModuleID` scalar type represents an identifier for an object of type CurrentModule.
type CurrentModuleID string

// The `DirectoryID` scalar type represents an identifier for an object of type Directory.
type DirectoryID string

// The `EnvVariableID` scalar type represents an identifier for an object of type EnvVariable.
type EnvVariableID string

// The `FieldTypeDefID` scalar type represents an identifier for an object of type FieldTypeDef.
type FieldTypeDefID string

// The `FileID` scalar type represents an identifier for an object of type File.
type FileID string

// The `FunctionArgID` scalar type represents an identifier for an object of type FunctionArg.
type FunctionArgID string

// The `FunctionCallArgValueID` scalar type represents an identifier for an object of type FunctionCallArgValue.
type FunctionCallArgValueID string

// The `FunctionCallID` scalar type represents an identifier for an object of type FunctionCall.
type FunctionCallID string

// The `FunctionID` scalar type represents an identifier for an object of type Function.
type FunctionID string

// The `GeneratedCodeID` scalar type represents an identifier for an object of type GeneratedCode.
type GeneratedCodeID string

// The `GitModuleSourceID` scalar type represents an identifier for an object of type GitModuleSource.
type GitModuleSourceID string

// The `GitRefID` scalar type represents an identifier for an object of type GitRef.
type GitRefID string

// The `GitRepositoryID` scalar type represents an identifier for an object of type GitRepository.
type GitRepositoryID string

// The `HostID` scalar type represents an identifier for an object of type Host.
type HostID string

// The `InputTypeDefID` scalar type represents an identifier for an object of type InputTypeDef.
type InputTypeDefID string

// The `InterfaceTypeDefID` scalar type represents an identifier for an object of type InterfaceTypeDef.
type InterfaceTypeDefID string

// An arbitrary JSON-encoded value.
type JSON string

// The `LabelID` scalar type represents an identifier for an object of type Label.
type LabelID string

// The `ListTypeDefID` scalar type represents an identifier for an object of type ListTypeDef.
type ListTypeDefID string

// The `LocalModuleSourceID` scalar type represents an identifier for an object of type LocalModuleSource.
type LocalModuleSourceID string

// The `ModuleDependencyID` scalar type represents an identifier for an object of type ModuleDependency.
type ModuleDependencyID string

// The `ModuleID` scalar type represents an identifier for an object of type Module.
type ModuleID string

// The `ModuleSourceID` scalar type represents an identifier for an object of type ModuleSource.
type ModuleSourceID string

// The `ObjectTypeDefID` scalar type represents an identifier for an object of type ObjectTypeDef.
type ObjectTypeDefID string

// The platform config OS and architecture in a Container.
//
// The format is [os]/[platform]/[version] (e.g., "darwin/arm64/v7", "windows/amd64", "linux/arm64").
type Platform string

// The `PortID` scalar type represents an identifier for an object of type Port.
type PortID string

// The `SecretID` scalar type represents an identifier for an object of type Secret.
type SecretID string

// The `ServiceID` scalar type represents an identifier for an object of type Service.
type ServiceID string

// The `SocketID` scalar type represents an identifier for an object of type Socket.
type SocketID string

// The `TerminalID` scalar type represents an identifier for an object of type Terminal.
type TerminalID string

// The `TypeDefID` scalar type represents an identifier for an object of type TypeDef.
type TypeDefID string

// The absence of a value.
//
// A Null Void is used as a placeholder for resolvers that do not return anything.
type Void string

// Key value object that represents a build argument.
type BuildArg struct {
	// The build argument name.
	Name string `json:"name"`

	// The build argument value.
	Value string `json:"value"`
}

// Key value object that represents a pipeline label.
type PipelineLabel struct {
	// Label name.
	Name string `json:"name"`

	// Label value.
	Value string `json:"value"`
}

// Port forwarding rules for tunneling network traffic.
type PortForward struct {
	// Destination port for traffic.
	Backend int `json:"backend"`

	// Port to expose to clients. If unspecified, a default will be chosen.
	Frontend int `json:"frontend"`

	// Transport layer protocol to use for traffic.
	Protocol NetworkProtocol `json:"protocol,omitempty"`
}

// A directory whose contents persist across runs.
type CacheVolume struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	id *CacheVolumeID
}

// A unique identifier for this CacheVolume.
func (r *CacheVolume) ID(ctx context.Context) (CacheVolumeID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response CacheVolumeID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *CacheVolume) XXX_GraphQLType() string {
	return "CacheVolume"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *CacheVolume) XXX_GraphQLIDType() string {
	return "CacheVolumeID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *CacheVolume) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *CacheVolume) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// An OCI-compatible container, also known as a Docker container.
type Container struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	envVariable *string
	export      *bool
	id          *ContainerID
	imageRef    *string
	label       *string
	platform    *Platform
	publish     *string
	stderr      *string
	stdout      *string
	sync        *ContainerID
	user        *string
	workdir     *string
}
type WithContainerFunc func(r *Container) *Container

// With calls the provided function with current Container.
//
// This is useful for reusability and readability by not breaking the calling chain.
func (r *Container) With(f WithContainerFunc) *Container {
	return f(r)
}

// Turn the container into a Service.
//
// Be sure to set any exposed ports before this conversion.
func (r *Container) AsService() *Service {
	q := r.Query.Select("asService")

	return &Service{
		Query:  q,
		Client: r.Client,
	}
}

// ContainerAsTarballOpts contains options for Container.AsTarball
type ContainerAsTarballOpts struct {
	// Identifiers for other platform specific containers.
	//
	// Used for multi-platform images.
	PlatformVariants []*Container
	// Force each layer of the image to use the specified compression algorithm.
	//
	// If this is unset, then if a layer already has a compressed blob in the engine's cache, that will be used (this can result in a mix of compression algorithms for different layers). If this is unset and a layer has no compressed blob in the engine's cache, then it will be compressed using Gzip.
	ForcedCompression ImageLayerCompression
	// Use the specified media types for the image's layers.
	//
	// Defaults to OCI, which is largely compatible with most recent container runtimes, but Docker may be needed for older runtimes without OCI support.
	MediaTypes ImageMediaTypes
}

// Returns a File representing the container serialized to a tarball.
func (r *Container) AsTarball(opts ...ContainerAsTarballOpts) *File {
	q := r.Query.Select("asTarball")
	for i := len(opts) - 1; i >= 0; i-- {
		// `platformVariants` optional argument
		if !querybuilder.IsZeroValue(opts[i].PlatformVariants) {
			q = q.Arg("platformVariants", opts[i].PlatformVariants)
		}
		// `forcedCompression` optional argument
		if !querybuilder.IsZeroValue(opts[i].ForcedCompression) {
			q = q.Arg("forcedCompression", opts[i].ForcedCompression)
		}
		// `mediaTypes` optional argument
		if !querybuilder.IsZeroValue(opts[i].MediaTypes) {
			q = q.Arg("mediaTypes", opts[i].MediaTypes)
		}
	}

	return &File{
		Query:  q,
		Client: r.Client,
	}
}

// ContainerBuildOpts contains options for Container.Build
type ContainerBuildOpts struct {
	// Path to the Dockerfile to use.
	Dockerfile string
	// Target build stage to build.
	Target string
	// Additional build arguments.
	BuildArgs []BuildArg
	// Secrets to pass to the build.
	//
	// They will be mounted at /run/secrets/[secret-name] in the build container
	//
	// They can be accessed in the Dockerfile using the "secret" mount type and mount path /run/secrets/[secret-name], e.g. RUN --mount=type=secret,id=my-secret curl http://example.com?token=$(cat /run/secrets/my-secret)
	Secrets []*Secret
}

// Initializes this container from a Dockerfile build.
func (r *Container) Build(context *Directory, opts ...ContainerBuildOpts) *Container {
	assertNotNil("context", context)
	q := r.Query.Select("build")
	for i := len(opts) - 1; i >= 0; i-- {
		// `dockerfile` optional argument
		if !querybuilder.IsZeroValue(opts[i].Dockerfile) {
			q = q.Arg("dockerfile", opts[i].Dockerfile)
		}
		// `target` optional argument
		if !querybuilder.IsZeroValue(opts[i].Target) {
			q = q.Arg("target", opts[i].Target)
		}
		// `buildArgs` optional argument
		if !querybuilder.IsZeroValue(opts[i].BuildArgs) {
			q = q.Arg("buildArgs", opts[i].BuildArgs)
		}
		// `secrets` optional argument
		if !querybuilder.IsZeroValue(opts[i].Secrets) {
			q = q.Arg("secrets", opts[i].Secrets)
		}
	}
	q = q.Arg("context", context)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves default arguments for future commands.
func (r *Container) DefaultArgs(ctx context.Context) ([]string, error) {
	q := r.Query.Select("defaultArgs")

	var response []string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Retrieves a directory at the given path.
//
// Mounts are included.
func (r *Container) Directory(path string) *Directory {
	q := r.Query.Select("directory")
	q = q.Arg("path", path)

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves entrypoint to be prepended to the arguments of all commands.
func (r *Container) Entrypoint(ctx context.Context) ([]string, error) {
	q := r.Query.Select("entrypoint")

	var response []string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Retrieves the value of the specified environment variable.
func (r *Container) EnvVariable(ctx context.Context, name string) (string, error) {
	if r.envVariable != nil {
		return *r.envVariable, nil
	}
	q := r.Query.Select("envVariable")
	q = q.Arg("name", name)

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Retrieves the list of environment variables passed to commands.
func (r *Container) EnvVariables(ctx context.Context) ([]EnvVariable, error) {
	q := r.Query.Select("envVariables")

	q = q.Select("id")

	type envVariables struct {
		Id EnvVariableID
	}

	convert := func(fields []envVariables) []EnvVariable {
		out := []EnvVariable{}

		for i := range fields {
			val := EnvVariable{id: &fields[i].Id}
			val.Query = querybuilder.Query().Select("loadEnvVariableFromID").Arg("id", fields[i].Id)
			val.Client = r.Client
			out = append(out, val)
		}

		return out
	}
	var response []envVariables

	q = q.Bind(&response)

	err := q.Execute(ctx, r.Client)
	if err != nil {
		return nil, err
	}

	return convert(response), nil
}

// EXPERIMENTAL API! Subject to change/removal at any time.
//
// Configures all available GPUs on the host to be accessible to this container.
//
// This currently works for Nvidia devices only.
func (r *Container) ExperimentalWithAllGPUs() *Container {
	q := r.Query.Select("experimentalWithAllGPUs")

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// EXPERIMENTAL API! Subject to change/removal at any time.
//
// Configures the provided list of devices to be accesible to this container.
//
// This currently works for Nvidia devices only.
func (r *Container) ExperimentalWithGPU(devices []string) *Container {
	q := r.Query.Select("experimentalWithGPU")
	q = q.Arg("devices", devices)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// ContainerExportOpts contains options for Container.Export
type ContainerExportOpts struct {
	// Identifiers for other platform specific containers.
	//
	// Used for multi-platform image.
	PlatformVariants []*Container
	// Force each layer of the exported image to use the specified compression algorithm.
	//
	// If this is unset, then if a layer already has a compressed blob in the engine's cache, that will be used (this can result in a mix of compression algorithms for different layers). If this is unset and a layer has no compressed blob in the engine's cache, then it will be compressed using Gzip.
	ForcedCompression ImageLayerCompression
	// Use the specified media types for the exported image's layers.
	//
	// Defaults to OCI, which is largely compatible with most recent container runtimes, but Docker may be needed for older runtimes without OCI support.
	MediaTypes ImageMediaTypes
}

// Writes the container as an OCI tarball to the destination file path on the host.
//
// Return true on success.
//
// It can also export platform variants.
func (r *Container) Export(ctx context.Context, path string, opts ...ContainerExportOpts) (bool, error) {
	if r.export != nil {
		return *r.export, nil
	}
	q := r.Query.Select("export")
	for i := len(opts) - 1; i >= 0; i-- {
		// `platformVariants` optional argument
		if !querybuilder.IsZeroValue(opts[i].PlatformVariants) {
			q = q.Arg("platformVariants", opts[i].PlatformVariants)
		}
		// `forcedCompression` optional argument
		if !querybuilder.IsZeroValue(opts[i].ForcedCompression) {
			q = q.Arg("forcedCompression", opts[i].ForcedCompression)
		}
		// `mediaTypes` optional argument
		if !querybuilder.IsZeroValue(opts[i].MediaTypes) {
			q = q.Arg("mediaTypes", opts[i].MediaTypes)
		}
	}
	q = q.Arg("path", path)

	var response bool

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Retrieves the list of exposed ports.
//
// This includes ports already exposed by the image, even if not explicitly added with dagger.
func (r *Container) ExposedPorts(ctx context.Context) ([]Port, error) {
	q := r.Query.Select("exposedPorts")

	q = q.Select("id")

	type exposedPorts struct {
		Id PortID
	}

	convert := func(fields []exposedPorts) []Port {
		out := []Port{}

		for i := range fields {
			val := Port{id: &fields[i].Id}
			val.Query = querybuilder.Query().Select("loadPortFromID").Arg("id", fields[i].Id)
			val.Client = r.Client
			out = append(out, val)
		}

		return out
	}
	var response []exposedPorts

	q = q.Bind(&response)

	err := q.Execute(ctx, r.Client)
	if err != nil {
		return nil, err
	}

	return convert(response), nil
}

// Retrieves a file at the given path.
//
// Mounts are included.
func (r *Container) File(path string) *File {
	q := r.Query.Select("file")
	q = q.Arg("path", path)

	return &File{
		Query:  q,
		Client: r.Client,
	}
}

// Initializes this container from a pulled base image.
func (r *Container) From(address string) *Container {
	q := r.Query.Select("from")
	q = q.Arg("address", address)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// A unique identifier for this Container.
func (r *Container) ID(ctx context.Context) (ContainerID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response ContainerID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *Container) XXX_GraphQLType() string {
	return "Container"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *Container) XXX_GraphQLIDType() string {
	return "ContainerID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *Container) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *Container) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The unique image reference which can only be retrieved immediately after the 'Container.From' call.
func (r *Container) ImageRef(ctx context.Context) (string, error) {
	if r.imageRef != nil {
		return *r.imageRef, nil
	}
	q := r.Query.Select("imageRef")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// ContainerImportOpts contains options for Container.Import
type ContainerImportOpts struct {
	// Identifies the tag to import from the archive, if the archive bundles multiple tags.
	Tag string
}

// Reads the container from an OCI tarball.
func (r *Container) Import(source *File, opts ...ContainerImportOpts) *Container {
	assertNotNil("source", source)
	q := r.Query.Select("import")
	for i := len(opts) - 1; i >= 0; i-- {
		// `tag` optional argument
		if !querybuilder.IsZeroValue(opts[i].Tag) {
			q = q.Arg("tag", opts[i].Tag)
		}
	}
	q = q.Arg("source", source)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves the value of the specified label.
func (r *Container) Label(ctx context.Context, name string) (string, error) {
	if r.label != nil {
		return *r.label, nil
	}
	q := r.Query.Select("label")
	q = q.Arg("name", name)

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Retrieves the list of labels passed to container.
func (r *Container) Labels(ctx context.Context) ([]Label, error) {
	q := r.Query.Select("labels")

	q = q.Select("id")

	type labels struct {
		Id LabelID
	}

	convert := func(fields []labels) []Label {
		out := []Label{}

		for i := range fields {
			val := Label{id: &fields[i].Id}
			val.Query = querybuilder.Query().Select("loadLabelFromID").Arg("id", fields[i].Id)
			val.Client = r.Client
			out = append(out, val)
		}

		return out
	}
	var response []labels

	q = q.Bind(&response)

	err := q.Execute(ctx, r.Client)
	if err != nil {
		return nil, err
	}

	return convert(response), nil
}

// Retrieves the list of paths where a directory is mounted.
func (r *Container) Mounts(ctx context.Context) ([]string, error) {
	q := r.Query.Select("mounts")

	var response []string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// ContainerPipelineOpts contains options for Container.Pipeline
type ContainerPipelineOpts struct {
	// Description of the sub-pipeline.
	Description string
	// Labels to apply to the sub-pipeline.
	Labels []PipelineLabel
}

// Creates a named sub-pipeline.
func (r *Container) Pipeline(name string, opts ...ContainerPipelineOpts) *Container {
	q := r.Query.Select("pipeline")
	for i := len(opts) - 1; i >= 0; i-- {
		// `description` optional argument
		if !querybuilder.IsZeroValue(opts[i].Description) {
			q = q.Arg("description", opts[i].Description)
		}
		// `labels` optional argument
		if !querybuilder.IsZeroValue(opts[i].Labels) {
			q = q.Arg("labels", opts[i].Labels)
		}
	}
	q = q.Arg("name", name)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// The platform this container executes and publishes as.
func (r *Container) Platform(ctx context.Context) (Platform, error) {
	if r.platform != nil {
		return *r.platform, nil
	}
	q := r.Query.Select("platform")

	var response Platform

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// ContainerPublishOpts contains options for Container.Publish
type ContainerPublishOpts struct {
	// Identifiers for other platform specific containers.
	//
	// Used for multi-platform image.
	PlatformVariants []*Container
	// Force each layer of the published image to use the specified compression algorithm.
	//
	// If this is unset, then if a layer already has a compressed blob in the engine's cache, that will be used (this can result in a mix of compression algorithms for different layers). If this is unset and a layer has no compressed blob in the engine's cache, then it will be compressed using Gzip.
	ForcedCompression ImageLayerCompression
	// Use the specified media types for the published image's layers.
	//
	// Defaults to OCI, which is largely compatible with most recent registries, but Docker may be needed for older registries without OCI support.
	MediaTypes ImageMediaTypes
}

// Publishes this container as a new image to the specified address.
//
// Publish returns a fully qualified ref.
//
// It can also publish platform variants.
func (r *Container) Publish(ctx context.Context, address string, opts ...ContainerPublishOpts) (string, error) {
	if r.publish != nil {
		return *r.publish, nil
	}
	q := r.Query.Select("publish")
	for i := len(opts) - 1; i >= 0; i-- {
		// `platformVariants` optional argument
		if !querybuilder.IsZeroValue(opts[i].PlatformVariants) {
			q = q.Arg("platformVariants", opts[i].PlatformVariants)
		}
		// `forcedCompression` optional argument
		if !querybuilder.IsZeroValue(opts[i].ForcedCompression) {
			q = q.Arg("forcedCompression", opts[i].ForcedCompression)
		}
		// `mediaTypes` optional argument
		if !querybuilder.IsZeroValue(opts[i].MediaTypes) {
			q = q.Arg("mediaTypes", opts[i].MediaTypes)
		}
	}
	q = q.Arg("address", address)

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Retrieves this container's root filesystem. Mounts are not included.
func (r *Container) Rootfs() *Directory {
	q := r.Query.Select("rootfs")

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// The error stream of the last executed command.
//
// Will execute default command if none is set, or error if there's no default.
func (r *Container) Stderr(ctx context.Context) (string, error) {
	if r.stderr != nil {
		return *r.stderr, nil
	}
	q := r.Query.Select("stderr")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The output stream of the last executed command.
//
// Will execute default command if none is set, or error if there's no default.
func (r *Container) Stdout(ctx context.Context) (string, error) {
	if r.stdout != nil {
		return *r.stdout, nil
	}
	q := r.Query.Select("stdout")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Forces evaluation of the pipeline in the engine.
//
// It doesn't run the default command if no exec has been set.
func (r *Container) Sync(ctx context.Context) (*Container, error) {
	q := r.Query.Select("sync")

	return r, q.Execute(ctx, r.Client)
}

// ContainerTerminalOpts contains options for Container.Terminal
type ContainerTerminalOpts struct {
	// If set, override the container's default terminal command and invoke these command arguments instead.
	Cmd []string
}

// Return an interactive terminal for this container using its configured default terminal command if not overridden by args (or sh as a fallback default).
func (r *Container) Terminal(opts ...ContainerTerminalOpts) *Terminal {
	q := r.Query.Select("terminal")
	for i := len(opts) - 1; i >= 0; i-- {
		// `cmd` optional argument
		if !querybuilder.IsZeroValue(opts[i].Cmd) {
			q = q.Arg("cmd", opts[i].Cmd)
		}
	}

	return &Terminal{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves the user to be set for all commands.
func (r *Container) User(ctx context.Context) (string, error) {
	if r.user != nil {
		return *r.user, nil
	}
	q := r.Query.Select("user")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Configures default arguments for future commands.
func (r *Container) WithDefaultArgs(args []string) *Container {
	q := r.Query.Select("withDefaultArgs")
	q = q.Arg("args", args)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Set the default command to invoke for the container's terminal API.
func (r *Container) WithDefaultTerminalCmd(args []string) *Container {
	q := r.Query.Select("withDefaultTerminalCmd")
	q = q.Arg("args", args)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// ContainerWithDirectoryOpts contains options for Container.WithDirectory
type ContainerWithDirectoryOpts struct {
	// Patterns to exclude in the written directory (e.g. ["node_modules/**", ".gitignore", ".git/"]).
	Exclude []string
	// Patterns to include in the written directory (e.g. ["*.go", "go.mod", "go.sum"]).
	Include []string
	// A user:group to set for the directory and its contents.
	//
	// The user and group can either be an ID (1000:1000) or a name (foo:bar).
	//
	// If the group is omitted, it defaults to the same as the user.
	Owner string
}

// Retrieves this container plus a directory written at the given path.
func (r *Container) WithDirectory(path string, directory *Directory, opts ...ContainerWithDirectoryOpts) *Container {
	assertNotNil("directory", directory)
	q := r.Query.Select("withDirectory")
	for i := len(opts) - 1; i >= 0; i-- {
		// `exclude` optional argument
		if !querybuilder.IsZeroValue(opts[i].Exclude) {
			q = q.Arg("exclude", opts[i].Exclude)
		}
		// `include` optional argument
		if !querybuilder.IsZeroValue(opts[i].Include) {
			q = q.Arg("include", opts[i].Include)
		}
		// `owner` optional argument
		if !querybuilder.IsZeroValue(opts[i].Owner) {
			q = q.Arg("owner", opts[i].Owner)
		}
	}
	q = q.Arg("path", path)
	q = q.Arg("directory", directory)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// ContainerWithEntrypointOpts contains options for Container.WithEntrypoint
type ContainerWithEntrypointOpts struct {
	// Don't remove the default arguments when setting the entrypoint.
	KeepDefaultArgs bool
}

// Retrieves this container but with a different command entrypoint.
func (r *Container) WithEntrypoint(args []string, opts ...ContainerWithEntrypointOpts) *Container {
	q := r.Query.Select("withEntrypoint")
	for i := len(opts) - 1; i >= 0; i-- {
		// `keepDefaultArgs` optional argument
		if !querybuilder.IsZeroValue(opts[i].KeepDefaultArgs) {
			q = q.Arg("keepDefaultArgs", opts[i].KeepDefaultArgs)
		}
	}
	q = q.Arg("args", args)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// ContainerWithEnvVariableOpts contains options for Container.WithEnvVariable
type ContainerWithEnvVariableOpts struct {
	// Replace `${VAR}` or `$VAR` in the value according to the current environment variables defined in the container (e.g., "/opt/bin:$PATH").
	Expand bool
}

// Retrieves this container plus the given environment variable.
func (r *Container) WithEnvVariable(name string, value string, opts ...ContainerWithEnvVariableOpts) *Container {
	q := r.Query.Select("withEnvVariable")
	for i := len(opts) - 1; i >= 0; i-- {
		// `expand` optional argument
		if !querybuilder.IsZeroValue(opts[i].Expand) {
			q = q.Arg("expand", opts[i].Expand)
		}
	}
	q = q.Arg("name", name)
	q = q.Arg("value", value)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// ContainerWithExecOpts contains options for Container.WithExec
type ContainerWithExecOpts struct {
	// If the container has an entrypoint, ignore it for args rather than using it to wrap them.
	SkipEntrypoint bool
	// Content to write to the command's standard input before closing (e.g., "Hello world").
	Stdin string
	// Redirect the command's standard output to a file in the container (e.g., "/tmp/stdout").
	RedirectStdout string
	// Redirect the command's standard error to a file in the container (e.g., "/tmp/stderr").
	RedirectStderr string
	// Provides dagger access to the executed command.
	//
	// Do not use this option unless you trust the command being executed; the command being executed WILL BE GRANTED FULL ACCESS TO YOUR HOST FILESYSTEM.
	ExperimentalPrivilegedNesting bool
	// Execute the command with all root capabilities. This is similar to running a command with "sudo" or executing "docker run" with the "--privileged" flag. Containerization does not provide any security guarantees when using this option. It should only be used when absolutely necessary and only with trusted commands.
	InsecureRootCapabilities bool
}

// Retrieves this container after executing the specified command inside it.
func (r *Container) WithExec(args []string, opts ...ContainerWithExecOpts) *Container {
	q := r.Query.Select("withExec")
	for i := len(opts) - 1; i >= 0; i-- {
		// `skipEntrypoint` optional argument
		if !querybuilder.IsZeroValue(opts[i].SkipEntrypoint) {
			q = q.Arg("skipEntrypoint", opts[i].SkipEntrypoint)
		}
		// `stdin` optional argument
		if !querybuilder.IsZeroValue(opts[i].Stdin) {
			q = q.Arg("stdin", opts[i].Stdin)
		}
		// `redirectStdout` optional argument
		if !querybuilder.IsZeroValue(opts[i].RedirectStdout) {
			q = q.Arg("redirectStdout", opts[i].RedirectStdout)
		}
		// `redirectStderr` optional argument
		if !querybuilder.IsZeroValue(opts[i].RedirectStderr) {
			q = q.Arg("redirectStderr", opts[i].RedirectStderr)
		}
		// `experimentalPrivilegedNesting` optional argument
		if !querybuilder.IsZeroValue(opts[i].ExperimentalPrivilegedNesting) {
			q = q.Arg("experimentalPrivilegedNesting", opts[i].ExperimentalPrivilegedNesting)
		}
		// `insecureRootCapabilities` optional argument
		if !querybuilder.IsZeroValue(opts[i].InsecureRootCapabilities) {
			q = q.Arg("insecureRootCapabilities", opts[i].InsecureRootCapabilities)
		}
	}
	q = q.Arg("args", args)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// ContainerWithExposedPortOpts contains options for Container.WithExposedPort
type ContainerWithExposedPortOpts struct {
	// Transport layer network protocol
	Protocol NetworkProtocol
	// Optional port description
	Description string
	// Skip the health check when run as a service.
	ExperimentalSkipHealthcheck bool
}

// Expose a network port.
//
// Exposed ports serve two purposes:
//
// - For health checks and introspection, when running services
//
// - For setting the EXPOSE OCI field when publishing the container
func (r *Container) WithExposedPort(port int, opts ...ContainerWithExposedPortOpts) *Container {
	q := r.Query.Select("withExposedPort")
	for i := len(opts) - 1; i >= 0; i-- {
		// `protocol` optional argument
		if !querybuilder.IsZeroValue(opts[i].Protocol) {
			q = q.Arg("protocol", opts[i].Protocol)
		}
		// `description` optional argument
		if !querybuilder.IsZeroValue(opts[i].Description) {
			q = q.Arg("description", opts[i].Description)
		}
		// `experimentalSkipHealthcheck` optional argument
		if !querybuilder.IsZeroValue(opts[i].ExperimentalSkipHealthcheck) {
			q = q.Arg("experimentalSkipHealthcheck", opts[i].ExperimentalSkipHealthcheck)
		}
	}
	q = q.Arg("port", port)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// ContainerWithFileOpts contains options for Container.WithFile
type ContainerWithFileOpts struct {
	// Permission given to the copied file (e.g., 0600).
	Permissions int
	// A user:group to set for the file.
	//
	// The user and group can either be an ID (1000:1000) or a name (foo:bar).
	//
	// If the group is omitted, it defaults to the same as the user.
	Owner string
}

// Retrieves this container plus the contents of the given file copied to the given path.
func (r *Container) WithFile(path string, source *File, opts ...ContainerWithFileOpts) *Container {
	assertNotNil("source", source)
	q := r.Query.Select("withFile")
	for i := len(opts) - 1; i >= 0; i-- {
		// `permissions` optional argument
		if !querybuilder.IsZeroValue(opts[i].Permissions) {
			q = q.Arg("permissions", opts[i].Permissions)
		}
		// `owner` optional argument
		if !querybuilder.IsZeroValue(opts[i].Owner) {
			q = q.Arg("owner", opts[i].Owner)
		}
	}
	q = q.Arg("path", path)
	q = q.Arg("source", source)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// ContainerWithFilesOpts contains options for Container.WithFiles
type ContainerWithFilesOpts struct {
	// Permission given to the copied files (e.g., 0600).
	Permissions int
	// A user:group to set for the files.
	//
	// The user and group can either be an ID (1000:1000) or a name (foo:bar).
	//
	// If the group is omitted, it defaults to the same as the user.
	Owner string
}

// Retrieves this container plus the contents of the given files copied to the given path.
func (r *Container) WithFiles(path string, sources []*File, opts ...ContainerWithFilesOpts) *Container {
	q := r.Query.Select("withFiles")
	for i := len(opts) - 1; i >= 0; i-- {
		// `permissions` optional argument
		if !querybuilder.IsZeroValue(opts[i].Permissions) {
			q = q.Arg("permissions", opts[i].Permissions)
		}
		// `owner` optional argument
		if !querybuilder.IsZeroValue(opts[i].Owner) {
			q = q.Arg("owner", opts[i].Owner)
		}
	}
	q = q.Arg("path", path)
	q = q.Arg("sources", sources)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Indicate that subsequent operations should be featured more prominently in the UI.
func (r *Container) WithFocus() *Container {
	q := r.Query.Select("withFocus")

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves this container plus the given label.
func (r *Container) WithLabel(name string, value string) *Container {
	q := r.Query.Select("withLabel")
	q = q.Arg("name", name)
	q = q.Arg("value", value)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// ContainerWithMountedCacheOpts contains options for Container.WithMountedCache
type ContainerWithMountedCacheOpts struct {
	// Identifier of the directory to use as the cache volume's root.
	Source *Directory
	// Sharing mode of the cache volume.
	Sharing CacheSharingMode
	// A user:group to set for the mounted cache directory.
	//
	// Note that this changes the ownership of the specified mount along with the initial filesystem provided by source (if any). It does not have any effect if/when the cache has already been created.
	//
	// The user and group can either be an ID (1000:1000) or a name (foo:bar).
	//
	// If the group is omitted, it defaults to the same as the user.
	Owner string
}

// Retrieves this container plus a cache volume mounted at the given path.
func (r *Container) WithMountedCache(path string, cache *CacheVolume, opts ...ContainerWithMountedCacheOpts) *Container {
	assertNotNil("cache", cache)
	q := r.Query.Select("withMountedCache")
	for i := len(opts) - 1; i >= 0; i-- {
		// `source` optional argument
		if !querybuilder.IsZeroValue(opts[i].Source) {
			q = q.Arg("source", opts[i].Source)
		}
		// `sharing` optional argument
		if !querybuilder.IsZeroValue(opts[i].Sharing) {
			q = q.Arg("sharing", opts[i].Sharing)
		}
		// `owner` optional argument
		if !querybuilder.IsZeroValue(opts[i].Owner) {
			q = q.Arg("owner", opts[i].Owner)
		}
	}
	q = q.Arg("path", path)
	q = q.Arg("cache", cache)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// ContainerWithMountedDirectoryOpts contains options for Container.WithMountedDirectory
type ContainerWithMountedDirectoryOpts struct {
	// A user:group to set for the mounted directory and its contents.
	//
	// The user and group can either be an ID (1000:1000) or a name (foo:bar).
	//
	// If the group is omitted, it defaults to the same as the user.
	Owner string
}

// Retrieves this container plus a directory mounted at the given path.
func (r *Container) WithMountedDirectory(path string, source *Directory, opts ...ContainerWithMountedDirectoryOpts) *Container {
	assertNotNil("source", source)
	q := r.Query.Select("withMountedDirectory")
	for i := len(opts) - 1; i >= 0; i-- {
		// `owner` optional argument
		if !querybuilder.IsZeroValue(opts[i].Owner) {
			q = q.Arg("owner", opts[i].Owner)
		}
	}
	q = q.Arg("path", path)
	q = q.Arg("source", source)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// ContainerWithMountedFileOpts contains options for Container.WithMountedFile
type ContainerWithMountedFileOpts struct {
	// A user or user:group to set for the mounted file.
	//
	// The user and group can either be an ID (1000:1000) or a name (foo:bar).
	//
	// If the group is omitted, it defaults to the same as the user.
	Owner string
}

// Retrieves this container plus a file mounted at the given path.
func (r *Container) WithMountedFile(path string, source *File, opts ...ContainerWithMountedFileOpts) *Container {
	assertNotNil("source", source)
	q := r.Query.Select("withMountedFile")
	for i := len(opts) - 1; i >= 0; i-- {
		// `owner` optional argument
		if !querybuilder.IsZeroValue(opts[i].Owner) {
			q = q.Arg("owner", opts[i].Owner)
		}
	}
	q = q.Arg("path", path)
	q = q.Arg("source", source)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// ContainerWithMountedSecretOpts contains options for Container.WithMountedSecret
type ContainerWithMountedSecretOpts struct {
	// A user:group to set for the mounted secret.
	//
	// The user and group can either be an ID (1000:1000) or a name (foo:bar).
	//
	// If the group is omitted, it defaults to the same as the user.
	Owner string
	// Permission given to the mounted secret (e.g., 0600).
	//
	// This option requires an owner to be set to be active.
	Mode int
}

// Retrieves this container plus a secret mounted into a file at the given path.
func (r *Container) WithMountedSecret(path string, source *Secret, opts ...ContainerWithMountedSecretOpts) *Container {
	assertNotNil("source", source)
	q := r.Query.Select("withMountedSecret")
	for i := len(opts) - 1; i >= 0; i-- {
		// `owner` optional argument
		if !querybuilder.IsZeroValue(opts[i].Owner) {
			q = q.Arg("owner", opts[i].Owner)
		}
		// `mode` optional argument
		if !querybuilder.IsZeroValue(opts[i].Mode) {
			q = q.Arg("mode", opts[i].Mode)
		}
	}
	q = q.Arg("path", path)
	q = q.Arg("source", source)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves this container plus a temporary directory mounted at the given path.
func (r *Container) WithMountedTemp(path string) *Container {
	q := r.Query.Select("withMountedTemp")
	q = q.Arg("path", path)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// ContainerWithNewFileOpts contains options for Container.WithNewFile
type ContainerWithNewFileOpts struct {
	// Content of the file to write (e.g., "Hello world!").
	Contents string
	// Permission given to the written file (e.g., 0600).
	Permissions int
	// A user:group to set for the file.
	//
	// The user and group can either be an ID (1000:1000) or a name (foo:bar).
	//
	// If the group is omitted, it defaults to the same as the user.
	Owner string
}

// Retrieves this container plus a new file written at the given path.
func (r *Container) WithNewFile(path string, opts ...ContainerWithNewFileOpts) *Container {
	q := r.Query.Select("withNewFile")
	for i := len(opts) - 1; i >= 0; i-- {
		// `contents` optional argument
		if !querybuilder.IsZeroValue(opts[i].Contents) {
			q = q.Arg("contents", opts[i].Contents)
		}
		// `permissions` optional argument
		if !querybuilder.IsZeroValue(opts[i].Permissions) {
			q = q.Arg("permissions", opts[i].Permissions)
		}
		// `owner` optional argument
		if !querybuilder.IsZeroValue(opts[i].Owner) {
			q = q.Arg("owner", opts[i].Owner)
		}
	}
	q = q.Arg("path", path)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves this container with a registry authentication for a given address.
func (r *Container) WithRegistryAuth(address string, username string, secret *Secret) *Container {
	assertNotNil("secret", secret)
	q := r.Query.Select("withRegistryAuth")
	q = q.Arg("address", address)
	q = q.Arg("username", username)
	q = q.Arg("secret", secret)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves the container with the given directory mounted to /.
func (r *Container) WithRootfs(directory *Directory) *Container {
	assertNotNil("directory", directory)
	q := r.Query.Select("withRootfs")
	q = q.Arg("directory", directory)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves this container plus an env variable containing the given secret.
func (r *Container) WithSecretVariable(name string, secret *Secret) *Container {
	assertNotNil("secret", secret)
	q := r.Query.Select("withSecretVariable")
	q = q.Arg("name", name)
	q = q.Arg("secret", secret)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Establish a runtime dependency on a service.
//
// The service will be started automatically when needed and detached when it is no longer needed, executing the default command if none is set.
//
// The service will be reachable from the container via the provided hostname alias.
//
// The service dependency will also convey to any files or directories produced by the container.
func (r *Container) WithServiceBinding(alias string, service *Service) *Container {
	assertNotNil("service", service)
	q := r.Query.Select("withServiceBinding")
	q = q.Arg("alias", alias)
	q = q.Arg("service", service)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// ContainerWithUnixSocketOpts contains options for Container.WithUnixSocket
type ContainerWithUnixSocketOpts struct {
	// A user:group to set for the mounted socket.
	//
	// The user and group can either be an ID (1000:1000) or a name (foo:bar).
	//
	// If the group is omitted, it defaults to the same as the user.
	Owner string
}

// Retrieves this container plus a socket forwarded to the given Unix socket path.
func (r *Container) WithUnixSocket(path string, source *Socket, opts ...ContainerWithUnixSocketOpts) *Container {
	assertNotNil("source", source)
	q := r.Query.Select("withUnixSocket")
	for i := len(opts) - 1; i >= 0; i-- {
		// `owner` optional argument
		if !querybuilder.IsZeroValue(opts[i].Owner) {
			q = q.Arg("owner", opts[i].Owner)
		}
	}
	q = q.Arg("path", path)
	q = q.Arg("source", source)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves this container with a different command user.
func (r *Container) WithUser(name string) *Container {
	q := r.Query.Select("withUser")
	q = q.Arg("name", name)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves this container with a different working directory.
func (r *Container) WithWorkdir(path string) *Container {
	q := r.Query.Select("withWorkdir")
	q = q.Arg("path", path)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves this container with unset default arguments for future commands.
func (r *Container) WithoutDefaultArgs() *Container {
	q := r.Query.Select("withoutDefaultArgs")

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// ContainerWithoutEntrypointOpts contains options for Container.WithoutEntrypoint
type ContainerWithoutEntrypointOpts struct {
	// Don't remove the default arguments when unsetting the entrypoint.
	KeepDefaultArgs bool
}

// Retrieves this container with an unset command entrypoint.
func (r *Container) WithoutEntrypoint(opts ...ContainerWithoutEntrypointOpts) *Container {
	q := r.Query.Select("withoutEntrypoint")
	for i := len(opts) - 1; i >= 0; i-- {
		// `keepDefaultArgs` optional argument
		if !querybuilder.IsZeroValue(opts[i].KeepDefaultArgs) {
			q = q.Arg("keepDefaultArgs", opts[i].KeepDefaultArgs)
		}
	}

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves this container minus the given environment variable.
func (r *Container) WithoutEnvVariable(name string) *Container {
	q := r.Query.Select("withoutEnvVariable")
	q = q.Arg("name", name)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// ContainerWithoutExposedPortOpts contains options for Container.WithoutExposedPort
type ContainerWithoutExposedPortOpts struct {
	// Port protocol to unexpose
	Protocol NetworkProtocol
}

// Unexpose a previously exposed port.
func (r *Container) WithoutExposedPort(port int, opts ...ContainerWithoutExposedPortOpts) *Container {
	q := r.Query.Select("withoutExposedPort")
	for i := len(opts) - 1; i >= 0; i-- {
		// `protocol` optional argument
		if !querybuilder.IsZeroValue(opts[i].Protocol) {
			q = q.Arg("protocol", opts[i].Protocol)
		}
	}
	q = q.Arg("port", port)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Indicate that subsequent operations should not be featured more prominently in the UI.
//
// This is the initial state of all containers.
func (r *Container) WithoutFocus() *Container {
	q := r.Query.Select("withoutFocus")

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves this container minus the given environment label.
func (r *Container) WithoutLabel(name string) *Container {
	q := r.Query.Select("withoutLabel")
	q = q.Arg("name", name)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves this container after unmounting everything at the given path.
func (r *Container) WithoutMount(path string) *Container {
	q := r.Query.Select("withoutMount")
	q = q.Arg("path", path)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves this container without the registry authentication of a given address.
func (r *Container) WithoutRegistryAuth(address string) *Container {
	q := r.Query.Select("withoutRegistryAuth")
	q = q.Arg("address", address)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves this container with a previously added Unix socket removed.
func (r *Container) WithoutUnixSocket(path string) *Container {
	q := r.Query.Select("withoutUnixSocket")
	q = q.Arg("path", path)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves this container with an unset command user.
//
// Should default to root.
func (r *Container) WithoutUser() *Container {
	q := r.Query.Select("withoutUser")

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves this container with an unset working directory.
//
// Should default to "/".
func (r *Container) WithoutWorkdir() *Container {
	q := r.Query.Select("withoutWorkdir")

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves the working directory for all commands.
func (r *Container) Workdir(ctx context.Context) (string, error) {
	if r.workdir != nil {
		return *r.workdir, nil
	}
	q := r.Query.Select("workdir")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Reflective module API provided to functions at runtime.
type CurrentModule struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	id   *CurrentModuleID
	name *string
}

// A unique identifier for this CurrentModule.
func (r *CurrentModule) ID(ctx context.Context) (CurrentModuleID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response CurrentModuleID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *CurrentModule) XXX_GraphQLType() string {
	return "CurrentModule"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *CurrentModule) XXX_GraphQLIDType() string {
	return "CurrentModuleID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *CurrentModule) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *CurrentModule) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The name of the module being executed in
func (r *CurrentModule) Name(ctx context.Context) (string, error) {
	if r.name != nil {
		return *r.name, nil
	}
	q := r.Query.Select("name")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The directory containing the module's source code loaded into the engine (plus any generated code that may have been created).
func (r *CurrentModule) Source() *Directory {
	q := r.Query.Select("source")

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// CurrentModuleWorkdirOpts contains options for CurrentModule.Workdir
type CurrentModuleWorkdirOpts struct {
	// Exclude artifacts that match the given pattern (e.g., ["node_modules/", ".git*"]).
	Exclude []string
	// Include only artifacts that match the given pattern (e.g., ["app/", "package.*"]).
	Include []string
}

// Load a directory from the module's scratch working directory, including any changes that may have been made to it during module function execution.
func (r *CurrentModule) Workdir(path string, opts ...CurrentModuleWorkdirOpts) *Directory {
	q := r.Query.Select("workdir")
	for i := len(opts) - 1; i >= 0; i-- {
		// `exclude` optional argument
		if !querybuilder.IsZeroValue(opts[i].Exclude) {
			q = q.Arg("exclude", opts[i].Exclude)
		}
		// `include` optional argument
		if !querybuilder.IsZeroValue(opts[i].Include) {
			q = q.Arg("include", opts[i].Include)
		}
	}
	q = q.Arg("path", path)

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// Load a file from the module's scratch working directory, including any changes that may have been made to it during module function execution.Load a file from the module's scratch working directory, including any changes that may have been made to it during module function execution.
func (r *CurrentModule) WorkdirFile(path string) *File {
	q := r.Query.Select("workdirFile")
	q = q.Arg("path", path)

	return &File{
		Query:  q,
		Client: r.Client,
	}
}

// A directory.
type Directory struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	export *bool
	id     *DirectoryID
	sync   *DirectoryID
}
type WithDirectoryFunc func(r *Directory) *Directory

// With calls the provided function with current Directory.
//
// This is useful for reusability and readability by not breaking the calling chain.
func (r *Directory) With(f WithDirectoryFunc) *Directory {
	return f(r)
}

// DirectoryAsModuleOpts contains options for Directory.AsModule
type DirectoryAsModuleOpts struct {
	// An optional subpath of the directory which contains the module's configuration file.
	//
	// This is needed when the module code is in a subdirectory but requires parent directories to be loaded in order to execute. For example, the module source code may need a go.mod, project.toml, package.json, etc. file from a parent directory.
	//
	// If not set, the module source code is loaded from the root of the directory.
	SourceRootPath string
}

// Load the directory as a Dagger module
func (r *Directory) AsModule(opts ...DirectoryAsModuleOpts) *Module {
	q := r.Query.Select("asModule")
	for i := len(opts) - 1; i >= 0; i-- {
		// `sourceRootPath` optional argument
		if !querybuilder.IsZeroValue(opts[i].SourceRootPath) {
			q = q.Arg("sourceRootPath", opts[i].SourceRootPath)
		}
	}

	return &Module{
		Query:  q,
		Client: r.Client,
	}
}

// Gets the difference between this directory and an another directory.
func (r *Directory) Diff(other *Directory) *Directory {
	assertNotNil("other", other)
	q := r.Query.Select("diff")
	q = q.Arg("other", other)

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves a directory at the given path.
func (r *Directory) Directory(path string) *Directory {
	q := r.Query.Select("directory")
	q = q.Arg("path", path)

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// DirectoryDockerBuildOpts contains options for Directory.DockerBuild
type DirectoryDockerBuildOpts struct {
	// The platform to build.
	Platform Platform
	// Path to the Dockerfile to use (e.g., "frontend.Dockerfile").
	Dockerfile string
	// Target build stage to build.
	Target string
	// Build arguments to use in the build.
	BuildArgs []BuildArg
	// Secrets to pass to the build.
	//
	// They will be mounted at /run/secrets/[secret-name].
	Secrets []*Secret
}

// Builds a new Docker container from this directory.
func (r *Directory) DockerBuild(opts ...DirectoryDockerBuildOpts) *Container {
	q := r.Query.Select("dockerBuild")
	for i := len(opts) - 1; i >= 0; i-- {
		// `platform` optional argument
		if !querybuilder.IsZeroValue(opts[i].Platform) {
			q = q.Arg("platform", opts[i].Platform)
		}
		// `dockerfile` optional argument
		if !querybuilder.IsZeroValue(opts[i].Dockerfile) {
			q = q.Arg("dockerfile", opts[i].Dockerfile)
		}
		// `target` optional argument
		if !querybuilder.IsZeroValue(opts[i].Target) {
			q = q.Arg("target", opts[i].Target)
		}
		// `buildArgs` optional argument
		if !querybuilder.IsZeroValue(opts[i].BuildArgs) {
			q = q.Arg("buildArgs", opts[i].BuildArgs)
		}
		// `secrets` optional argument
		if !querybuilder.IsZeroValue(opts[i].Secrets) {
			q = q.Arg("secrets", opts[i].Secrets)
		}
	}

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// DirectoryEntriesOpts contains options for Directory.Entries
type DirectoryEntriesOpts struct {
	// Location of the directory to look at (e.g., "/src").
	Path string
}

// Returns a list of files and directories at the given path.
func (r *Directory) Entries(ctx context.Context, opts ...DirectoryEntriesOpts) ([]string, error) {
	q := r.Query.Select("entries")
	for i := len(opts) - 1; i >= 0; i-- {
		// `path` optional argument
		if !querybuilder.IsZeroValue(opts[i].Path) {
			q = q.Arg("path", opts[i].Path)
		}
	}

	var response []string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Writes the contents of the directory to a path on the host.
func (r *Directory) Export(ctx context.Context, path string) (bool, error) {
	if r.export != nil {
		return *r.export, nil
	}
	q := r.Query.Select("export")
	q = q.Arg("path", path)

	var response bool

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Retrieves a file at the given path.
func (r *Directory) File(path string) *File {
	q := r.Query.Select("file")
	q = q.Arg("path", path)

	return &File{
		Query:  q,
		Client: r.Client,
	}
}

// Returns a list of files and directories that matche the given pattern.
func (r *Directory) Glob(ctx context.Context, pattern string) ([]string, error) {
	q := r.Query.Select("glob")
	q = q.Arg("pattern", pattern)

	var response []string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A unique identifier for this Directory.
func (r *Directory) ID(ctx context.Context) (DirectoryID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response DirectoryID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *Directory) XXX_GraphQLType() string {
	return "Directory"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *Directory) XXX_GraphQLIDType() string {
	return "DirectoryID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *Directory) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *Directory) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// DirectoryPipelineOpts contains options for Directory.Pipeline
type DirectoryPipelineOpts struct {
	// Description of the sub-pipeline.
	Description string
	// Labels to apply to the sub-pipeline.
	Labels []PipelineLabel
}

// Creates a named sub-pipeline.
func (r *Directory) Pipeline(name string, opts ...DirectoryPipelineOpts) *Directory {
	q := r.Query.Select("pipeline")
	for i := len(opts) - 1; i >= 0; i-- {
		// `description` optional argument
		if !querybuilder.IsZeroValue(opts[i].Description) {
			q = q.Arg("description", opts[i].Description)
		}
		// `labels` optional argument
		if !querybuilder.IsZeroValue(opts[i].Labels) {
			q = q.Arg("labels", opts[i].Labels)
		}
	}
	q = q.Arg("name", name)

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// Force evaluation in the engine.
func (r *Directory) Sync(ctx context.Context) (*Directory, error) {
	q := r.Query.Select("sync")

	return r, q.Execute(ctx, r.Client)
}

// DirectoryWithDirectoryOpts contains options for Directory.WithDirectory
type DirectoryWithDirectoryOpts struct {
	// Exclude artifacts that match the given pattern (e.g., ["node_modules/", ".git*"]).
	Exclude []string
	// Include only artifacts that match the given pattern (e.g., ["app/", "package.*"]).
	Include []string
}

// Retrieves this directory plus a directory written at the given path.
func (r *Directory) WithDirectory(path string, directory *Directory, opts ...DirectoryWithDirectoryOpts) *Directory {
	assertNotNil("directory", directory)
	q := r.Query.Select("withDirectory")
	for i := len(opts) - 1; i >= 0; i-- {
		// `exclude` optional argument
		if !querybuilder.IsZeroValue(opts[i].Exclude) {
			q = q.Arg("exclude", opts[i].Exclude)
		}
		// `include` optional argument
		if !querybuilder.IsZeroValue(opts[i].Include) {
			q = q.Arg("include", opts[i].Include)
		}
	}
	q = q.Arg("path", path)
	q = q.Arg("directory", directory)

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// DirectoryWithFileOpts contains options for Directory.WithFile
type DirectoryWithFileOpts struct {
	// Permission given to the copied file (e.g., 0600).
	Permissions int
}

// Retrieves this directory plus the contents of the given file copied to the given path.
func (r *Directory) WithFile(path string, source *File, opts ...DirectoryWithFileOpts) *Directory {
	assertNotNil("source", source)
	q := r.Query.Select("withFile")
	for i := len(opts) - 1; i >= 0; i-- {
		// `permissions` optional argument
		if !querybuilder.IsZeroValue(opts[i].Permissions) {
			q = q.Arg("permissions", opts[i].Permissions)
		}
	}
	q = q.Arg("path", path)
	q = q.Arg("source", source)

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// DirectoryWithFilesOpts contains options for Directory.WithFiles
type DirectoryWithFilesOpts struct {
	// Permission given to the copied files (e.g., 0600).
	Permissions int
}

// Retrieves this directory plus the contents of the given files copied to the given path.
func (r *Directory) WithFiles(path string, sources []*File, opts ...DirectoryWithFilesOpts) *Directory {
	q := r.Query.Select("withFiles")
	for i := len(opts) - 1; i >= 0; i-- {
		// `permissions` optional argument
		if !querybuilder.IsZeroValue(opts[i].Permissions) {
			q = q.Arg("permissions", opts[i].Permissions)
		}
	}
	q = q.Arg("path", path)
	q = q.Arg("sources", sources)

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// DirectoryWithNewDirectoryOpts contains options for Directory.WithNewDirectory
type DirectoryWithNewDirectoryOpts struct {
	// Permission granted to the created directory (e.g., 0777).
	Permissions int
}

// Retrieves this directory plus a new directory created at the given path.
func (r *Directory) WithNewDirectory(path string, opts ...DirectoryWithNewDirectoryOpts) *Directory {
	q := r.Query.Select("withNewDirectory")
	for i := len(opts) - 1; i >= 0; i-- {
		// `permissions` optional argument
		if !querybuilder.IsZeroValue(opts[i].Permissions) {
			q = q.Arg("permissions", opts[i].Permissions)
		}
	}
	q = q.Arg("path", path)

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// DirectoryWithNewFileOpts contains options for Directory.WithNewFile
type DirectoryWithNewFileOpts struct {
	// Permission given to the copied file (e.g., 0600).
	Permissions int
}

// Retrieves this directory plus a new file written at the given path.
func (r *Directory) WithNewFile(path string, contents string, opts ...DirectoryWithNewFileOpts) *Directory {
	q := r.Query.Select("withNewFile")
	for i := len(opts) - 1; i >= 0; i-- {
		// `permissions` optional argument
		if !querybuilder.IsZeroValue(opts[i].Permissions) {
			q = q.Arg("permissions", opts[i].Permissions)
		}
	}
	q = q.Arg("path", path)
	q = q.Arg("contents", contents)

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves this directory with all file/dir timestamps set to the given time.
func (r *Directory) WithTimestamps(timestamp int) *Directory {
	q := r.Query.Select("withTimestamps")
	q = q.Arg("timestamp", timestamp)

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves this directory with the directory at the given path removed.
func (r *Directory) WithoutDirectory(path string) *Directory {
	q := r.Query.Select("withoutDirectory")
	q = q.Arg("path", path)

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves this directory with the file at the given path removed.
func (r *Directory) WithoutFile(path string) *Directory {
	q := r.Query.Select("withoutFile")
	q = q.Arg("path", path)

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// An environment variable name and value.
type EnvVariable struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	id    *EnvVariableID
	name  *string
	value *string
}

// A unique identifier for this EnvVariable.
func (r *EnvVariable) ID(ctx context.Context) (EnvVariableID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response EnvVariableID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *EnvVariable) XXX_GraphQLType() string {
	return "EnvVariable"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *EnvVariable) XXX_GraphQLIDType() string {
	return "EnvVariableID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *EnvVariable) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *EnvVariable) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The environment variable name.
func (r *EnvVariable) Name(ctx context.Context) (string, error) {
	if r.name != nil {
		return *r.name, nil
	}
	q := r.Query.Select("name")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The environment variable value.
func (r *EnvVariable) Value(ctx context.Context) (string, error) {
	if r.value != nil {
		return *r.value, nil
	}
	q := r.Query.Select("value")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A definition of a field on a custom object defined in a Module.
//
// A field on an object has a static value, as opposed to a function on an object whose value is computed by invoking code (and can accept arguments).
type FieldTypeDef struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	description *string
	id          *FieldTypeDefID
	name        *string
}

// A doc string for the field, if any.
func (r *FieldTypeDef) Description(ctx context.Context) (string, error) {
	if r.description != nil {
		return *r.description, nil
	}
	q := r.Query.Select("description")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A unique identifier for this FieldTypeDef.
func (r *FieldTypeDef) ID(ctx context.Context) (FieldTypeDefID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response FieldTypeDefID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *FieldTypeDef) XXX_GraphQLType() string {
	return "FieldTypeDef"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *FieldTypeDef) XXX_GraphQLIDType() string {
	return "FieldTypeDefID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *FieldTypeDef) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *FieldTypeDef) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The name of the field in lowerCamelCase format.
func (r *FieldTypeDef) Name(ctx context.Context) (string, error) {
	if r.name != nil {
		return *r.name, nil
	}
	q := r.Query.Select("name")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The type of the field.
func (r *FieldTypeDef) TypeDef() *TypeDef {
	q := r.Query.Select("typeDef")

	return &TypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// A file.
type File struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	contents *string
	export   *bool
	id       *FileID
	name     *string
	size     *int
	sync     *FileID
}
type WithFileFunc func(r *File) *File

// With calls the provided function with current File.
//
// This is useful for reusability and readability by not breaking the calling chain.
func (r *File) With(f WithFileFunc) *File {
	return f(r)
}

// Retrieves the contents of the file.
func (r *File) Contents(ctx context.Context) (string, error) {
	if r.contents != nil {
		return *r.contents, nil
	}
	q := r.Query.Select("contents")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// FileExportOpts contains options for File.Export
type FileExportOpts struct {
	// If allowParentDirPath is true, the path argument can be a directory path, in which case the file will be created in that directory.
	AllowParentDirPath bool
}

// Writes the file to a file path on the host.
func (r *File) Export(ctx context.Context, path string, opts ...FileExportOpts) (bool, error) {
	if r.export != nil {
		return *r.export, nil
	}
	q := r.Query.Select("export")
	for i := len(opts) - 1; i >= 0; i-- {
		// `allowParentDirPath` optional argument
		if !querybuilder.IsZeroValue(opts[i].AllowParentDirPath) {
			q = q.Arg("allowParentDirPath", opts[i].AllowParentDirPath)
		}
	}
	q = q.Arg("path", path)

	var response bool

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A unique identifier for this File.
func (r *File) ID(ctx context.Context) (FileID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response FileID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *File) XXX_GraphQLType() string {
	return "File"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *File) XXX_GraphQLIDType() string {
	return "FileID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *File) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *File) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// Retrieves the name of the file.
func (r *File) Name(ctx context.Context) (string, error) {
	if r.name != nil {
		return *r.name, nil
	}
	q := r.Query.Select("name")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Retrieves the size of the file, in bytes.
func (r *File) Size(ctx context.Context) (int, error) {
	if r.size != nil {
		return *r.size, nil
	}
	q := r.Query.Select("size")

	var response int

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Force evaluation in the engine.
func (r *File) Sync(ctx context.Context) (*File, error) {
	q := r.Query.Select("sync")

	return r, q.Execute(ctx, r.Client)
}

// Retrieves this file with its created/modified timestamps set to the given time.
func (r *File) WithTimestamps(timestamp int) *File {
	q := r.Query.Select("withTimestamps")
	q = q.Arg("timestamp", timestamp)

	return &File{
		Query:  q,
		Client: r.Client,
	}
}

// Function represents a resolver provided by a Module.
//
// A function always evaluates against a parent object and is given a set of named arguments.
type Function struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	description *string
	id          *FunctionID
	name        *string
}
type WithFunctionFunc func(r *Function) *Function

// With calls the provided function with current Function.
//
// This is useful for reusability and readability by not breaking the calling chain.
func (r *Function) With(f WithFunctionFunc) *Function {
	return f(r)
}

// Arguments accepted by the function, if any.
func (r *Function) Args(ctx context.Context) ([]FunctionArg, error) {
	q := r.Query.Select("args")

	q = q.Select("id")

	type args struct {
		Id FunctionArgID
	}

	convert := func(fields []args) []FunctionArg {
		out := []FunctionArg{}

		for i := range fields {
			val := FunctionArg{id: &fields[i].Id}
			val.Query = querybuilder.Query().Select("loadFunctionArgFromID").Arg("id", fields[i].Id)
			val.Client = r.Client
			out = append(out, val)
		}

		return out
	}
	var response []args

	q = q.Bind(&response)

	err := q.Execute(ctx, r.Client)
	if err != nil {
		return nil, err
	}

	return convert(response), nil
}

// A doc string for the function, if any.
func (r *Function) Description(ctx context.Context) (string, error) {
	if r.description != nil {
		return *r.description, nil
	}
	q := r.Query.Select("description")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A unique identifier for this Function.
func (r *Function) ID(ctx context.Context) (FunctionID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response FunctionID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *Function) XXX_GraphQLType() string {
	return "Function"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *Function) XXX_GraphQLIDType() string {
	return "FunctionID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *Function) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *Function) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The name of the function.
func (r *Function) Name(ctx context.Context) (string, error) {
	if r.name != nil {
		return *r.name, nil
	}
	q := r.Query.Select("name")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The type returned by the function.
func (r *Function) ReturnType() *TypeDef {
	q := r.Query.Select("returnType")

	return &TypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// FunctionWithArgOpts contains options for Function.WithArg
type FunctionWithArgOpts struct {
	// A doc string for the argument, if any
	Description string
	// A default value to use for this argument if not explicitly set by the caller, if any
	DefaultValue JSON
}

// Returns the function with the provided argument
func (r *Function) WithArg(name string, typeDef *TypeDef, opts ...FunctionWithArgOpts) *Function {
	assertNotNil("typeDef", typeDef)
	q := r.Query.Select("withArg")
	for i := len(opts) - 1; i >= 0; i-- {
		// `description` optional argument
		if !querybuilder.IsZeroValue(opts[i].Description) {
			q = q.Arg("description", opts[i].Description)
		}
		// `defaultValue` optional argument
		if !querybuilder.IsZeroValue(opts[i].DefaultValue) {
			q = q.Arg("defaultValue", opts[i].DefaultValue)
		}
	}
	q = q.Arg("name", name)
	q = q.Arg("typeDef", typeDef)

	return &Function{
		Query:  q,
		Client: r.Client,
	}
}

// Returns the function with the given doc string.
func (r *Function) WithDescription(description string) *Function {
	q := r.Query.Select("withDescription")
	q = q.Arg("description", description)

	return &Function{
		Query:  q,
		Client: r.Client,
	}
}

// An argument accepted by a function.
//
// This is a specification for an argument at function definition time, not an argument passed at function call time.
type FunctionArg struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	defaultValue *JSON
	description  *string
	id           *FunctionArgID
	name         *string
}

// A default value to use for this argument when not explicitly set by the caller, if any.
func (r *FunctionArg) DefaultValue(ctx context.Context) (JSON, error) {
	if r.defaultValue != nil {
		return *r.defaultValue, nil
	}
	q := r.Query.Select("defaultValue")

	var response JSON

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A doc string for the argument, if any.
func (r *FunctionArg) Description(ctx context.Context) (string, error) {
	if r.description != nil {
		return *r.description, nil
	}
	q := r.Query.Select("description")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A unique identifier for this FunctionArg.
func (r *FunctionArg) ID(ctx context.Context) (FunctionArgID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response FunctionArgID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *FunctionArg) XXX_GraphQLType() string {
	return "FunctionArg"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *FunctionArg) XXX_GraphQLIDType() string {
	return "FunctionArgID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *FunctionArg) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *FunctionArg) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The name of the argument in lowerCamelCase format.
func (r *FunctionArg) Name(ctx context.Context) (string, error) {
	if r.name != nil {
		return *r.name, nil
	}
	q := r.Query.Select("name")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The type of the argument.
func (r *FunctionArg) TypeDef() *TypeDef {
	q := r.Query.Select("typeDef")

	return &TypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// An active function call.
type FunctionCall struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	id          *FunctionCallID
	name        *string
	parent      *JSON
	parentName  *string
	returnValue *Void
}

// A unique identifier for this FunctionCall.
func (r *FunctionCall) ID(ctx context.Context) (FunctionCallID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response FunctionCallID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *FunctionCall) XXX_GraphQLType() string {
	return "FunctionCall"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *FunctionCall) XXX_GraphQLIDType() string {
	return "FunctionCallID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *FunctionCall) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *FunctionCall) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The argument values the function is being invoked with.
func (r *FunctionCall) InputArgs(ctx context.Context) ([]FunctionCallArgValue, error) {
	q := r.Query.Select("inputArgs")

	q = q.Select("id")

	type inputArgs struct {
		Id FunctionCallArgValueID
	}

	convert := func(fields []inputArgs) []FunctionCallArgValue {
		out := []FunctionCallArgValue{}

		for i := range fields {
			val := FunctionCallArgValue{id: &fields[i].Id}
			val.Query = querybuilder.Query().Select("loadFunctionCallArgValueFromID").Arg("id", fields[i].Id)
			val.Client = r.Client
			out = append(out, val)
		}

		return out
	}
	var response []inputArgs

	q = q.Bind(&response)

	err := q.Execute(ctx, r.Client)
	if err != nil {
		return nil, err
	}

	return convert(response), nil
}

// The name of the function being called.
func (r *FunctionCall) Name(ctx context.Context) (string, error) {
	if r.name != nil {
		return *r.name, nil
	}
	q := r.Query.Select("name")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The value of the parent object of the function being called. If the function is top-level to the module, this is always an empty object.
func (r *FunctionCall) Parent(ctx context.Context) (JSON, error) {
	if r.parent != nil {
		return *r.parent, nil
	}
	q := r.Query.Select("parent")

	var response JSON

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The name of the parent object of the function being called. If the function is top-level to the module, this is the name of the module.
func (r *FunctionCall) ParentName(ctx context.Context) (string, error) {
	if r.parentName != nil {
		return *r.parentName, nil
	}
	q := r.Query.Select("parentName")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Set the return value of the function call to the provided value.
func (r *FunctionCall) ReturnValue(ctx context.Context, value JSON) (Void, error) {
	if r.returnValue != nil {
		return *r.returnValue, nil
	}
	q := r.Query.Select("returnValue")
	q = q.Arg("value", value)

	var response Void

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A value passed as a named argument to a function call.
type FunctionCallArgValue struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	id    *FunctionCallArgValueID
	name  *string
	value *JSON
}

// A unique identifier for this FunctionCallArgValue.
func (r *FunctionCallArgValue) ID(ctx context.Context) (FunctionCallArgValueID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response FunctionCallArgValueID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *FunctionCallArgValue) XXX_GraphQLType() string {
	return "FunctionCallArgValue"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *FunctionCallArgValue) XXX_GraphQLIDType() string {
	return "FunctionCallArgValueID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *FunctionCallArgValue) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *FunctionCallArgValue) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The name of the argument.
func (r *FunctionCallArgValue) Name(ctx context.Context) (string, error) {
	if r.name != nil {
		return *r.name, nil
	}
	q := r.Query.Select("name")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The value of the argument represented as a JSON serialized string.
func (r *FunctionCallArgValue) Value(ctx context.Context) (JSON, error) {
	if r.value != nil {
		return *r.value, nil
	}
	q := r.Query.Select("value")

	var response JSON

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The result of running an SDK's codegen.
type GeneratedCode struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	id *GeneratedCodeID
}
type WithGeneratedCodeFunc func(r *GeneratedCode) *GeneratedCode

// With calls the provided function with current GeneratedCode.
//
// This is useful for reusability and readability by not breaking the calling chain.
func (r *GeneratedCode) With(f WithGeneratedCodeFunc) *GeneratedCode {
	return f(r)
}

// The directory containing the generated code.
func (r *GeneratedCode) Code() *Directory {
	q := r.Query.Select("code")

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// A unique identifier for this GeneratedCode.
func (r *GeneratedCode) ID(ctx context.Context) (GeneratedCodeID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response GeneratedCodeID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *GeneratedCode) XXX_GraphQLType() string {
	return "GeneratedCode"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *GeneratedCode) XXX_GraphQLIDType() string {
	return "GeneratedCodeID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *GeneratedCode) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *GeneratedCode) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// List of paths to mark generated in version control (i.e. .gitattributes).
func (r *GeneratedCode) VcsGeneratedPaths(ctx context.Context) ([]string, error) {
	q := r.Query.Select("vcsGeneratedPaths")

	var response []string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// List of paths to ignore in version control (i.e. .gitignore).
func (r *GeneratedCode) VcsIgnoredPaths(ctx context.Context) ([]string, error) {
	q := r.Query.Select("vcsIgnoredPaths")

	var response []string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Set the list of paths to mark generated in version control.
func (r *GeneratedCode) WithVCSGeneratedPaths(paths []string) *GeneratedCode {
	q := r.Query.Select("withVCSGeneratedPaths")
	q = q.Arg("paths", paths)

	return &GeneratedCode{
		Query:  q,
		Client: r.Client,
	}
}

// Set the list of paths to ignore in version control.
func (r *GeneratedCode) WithVCSIgnoredPaths(paths []string) *GeneratedCode {
	q := r.Query.Select("withVCSIgnoredPaths")
	q = q.Arg("paths", paths)

	return &GeneratedCode{
		Query:  q,
		Client: r.Client,
	}
}

// Module source originating from a git repo.
type GitModuleSource struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	cloneURL    *string
	commit      *string
	htmlURL     *string
	id          *GitModuleSourceID
	rootSubpath *string
	version     *string
}

// The URL from which the source's git repo can be cloned.
func (r *GitModuleSource) CloneURL(ctx context.Context) (string, error) {
	if r.cloneURL != nil {
		return *r.cloneURL, nil
	}
	q := r.Query.Select("cloneURL")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The resolved commit of the git repo this source points to.
func (r *GitModuleSource) Commit(ctx context.Context) (string, error) {
	if r.commit != nil {
		return *r.commit, nil
	}
	q := r.Query.Select("commit")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The directory containing everything needed to load load and use the module.
func (r *GitModuleSource) ContextDirectory() *Directory {
	q := r.Query.Select("contextDirectory")

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// The URL to the source's git repo in a web browser
func (r *GitModuleSource) HTMLURL(ctx context.Context) (string, error) {
	if r.htmlURL != nil {
		return *r.htmlURL, nil
	}
	q := r.Query.Select("htmlURL")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A unique identifier for this GitModuleSource.
func (r *GitModuleSource) ID(ctx context.Context) (GitModuleSourceID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response GitModuleSourceID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *GitModuleSource) XXX_GraphQLType() string {
	return "GitModuleSource"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *GitModuleSource) XXX_GraphQLIDType() string {
	return "GitModuleSourceID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *GitModuleSource) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *GitModuleSource) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The path to the root of the module source under the context directory. This directory contains its configuration file. It also contains its source code (possibly as a subdirectory).
func (r *GitModuleSource) RootSubpath(ctx context.Context) (string, error) {
	if r.rootSubpath != nil {
		return *r.rootSubpath, nil
	}
	q := r.Query.Select("rootSubpath")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The specified version of the git repo this source points to.
func (r *GitModuleSource) Version(ctx context.Context) (string, error) {
	if r.version != nil {
		return *r.version, nil
	}
	q := r.Query.Select("version")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A git ref (tag, branch, or commit).
type GitRef struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	commit *string
	id     *GitRefID
}

// The resolved commit id at this ref.
func (r *GitRef) Commit(ctx context.Context) (string, error) {
	if r.commit != nil {
		return *r.commit, nil
	}
	q := r.Query.Select("commit")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A unique identifier for this GitRef.
func (r *GitRef) ID(ctx context.Context) (GitRefID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response GitRefID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *GitRef) XXX_GraphQLType() string {
	return "GitRef"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *GitRef) XXX_GraphQLIDType() string {
	return "GitRefID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *GitRef) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *GitRef) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// GitRefTreeOpts contains options for GitRef.Tree
type GitRefTreeOpts struct {
	// DEPRECATED: This option should be passed to `git` instead.
	SSHKnownHosts string
	// DEPRECATED: This option should be passed to `git` instead.
	SSHAuthSocket *Socket
}

// The filesystem tree at this ref.
func (r *GitRef) Tree(opts ...GitRefTreeOpts) *Directory {
	q := r.Query.Select("tree")
	for i := len(opts) - 1; i >= 0; i-- {
		// `sshKnownHosts` optional argument
		if !querybuilder.IsZeroValue(opts[i].SSHKnownHosts) {
			q = q.Arg("sshKnownHosts", opts[i].SSHKnownHosts)
		}
		// `sshAuthSocket` optional argument
		if !querybuilder.IsZeroValue(opts[i].SSHAuthSocket) {
			q = q.Arg("sshAuthSocket", opts[i].SSHAuthSocket)
		}
	}

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// A git repository.
type GitRepository struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	id *GitRepositoryID
}

// Returns details of a branch.
func (r *GitRepository) Branch(name string) *GitRef {
	q := r.Query.Select("branch")
	q = q.Arg("name", name)

	return &GitRef{
		Query:  q,
		Client: r.Client,
	}
}

// Returns details of a commit.
func (r *GitRepository) Commit(id string) *GitRef {
	q := r.Query.Select("commit")
	q = q.Arg("id", id)

	return &GitRef{
		Query:  q,
		Client: r.Client,
	}
}

// A unique identifier for this GitRepository.
func (r *GitRepository) ID(ctx context.Context) (GitRepositoryID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response GitRepositoryID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *GitRepository) XXX_GraphQLType() string {
	return "GitRepository"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *GitRepository) XXX_GraphQLIDType() string {
	return "GitRepositoryID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *GitRepository) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *GitRepository) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// Returns details of a ref.
func (r *GitRepository) Ref(name string) *GitRef {
	q := r.Query.Select("ref")
	q = q.Arg("name", name)

	return &GitRef{
		Query:  q,
		Client: r.Client,
	}
}

// Returns details of a tag.
func (r *GitRepository) Tag(name string) *GitRef {
	q := r.Query.Select("tag")
	q = q.Arg("name", name)

	return &GitRef{
		Query:  q,
		Client: r.Client,
	}
}

// Information about the host environment.
type Host struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	id *HostID
}

// HostDirectoryOpts contains options for Host.Directory
type HostDirectoryOpts struct {
	// Exclude artifacts that match the given pattern (e.g., ["node_modules/", ".git*"]).
	Exclude []string
	// Include only artifacts that match the given pattern (e.g., ["app/", "package.*"]).
	Include []string
}

// Accesses a directory on the host.
func (r *Host) Directory(path string, opts ...HostDirectoryOpts) *Directory {
	q := r.Query.Select("directory")
	for i := len(opts) - 1; i >= 0; i-- {
		// `exclude` optional argument
		if !querybuilder.IsZeroValue(opts[i].Exclude) {
			q = q.Arg("exclude", opts[i].Exclude)
		}
		// `include` optional argument
		if !querybuilder.IsZeroValue(opts[i].Include) {
			q = q.Arg("include", opts[i].Include)
		}
	}
	q = q.Arg("path", path)

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// Accesses a file on the host.
func (r *Host) File(path string) *File {
	q := r.Query.Select("file")
	q = q.Arg("path", path)

	return &File{
		Query:  q,
		Client: r.Client,
	}
}

// A unique identifier for this Host.
func (r *Host) ID(ctx context.Context) (HostID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response HostID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *Host) XXX_GraphQLType() string {
	return "Host"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *Host) XXX_GraphQLIDType() string {
	return "HostID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *Host) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *Host) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// HostServiceOpts contains options for Host.Service
type HostServiceOpts struct {
	// Upstream host to forward traffic to.
	Host string
}

// Creates a service that forwards traffic to a specified address via the host.
func (r *Host) Service(ports []PortForward, opts ...HostServiceOpts) *Service {
	q := r.Query.Select("service")
	for i := len(opts) - 1; i >= 0; i-- {
		// `host` optional argument
		if !querybuilder.IsZeroValue(opts[i].Host) {
			q = q.Arg("host", opts[i].Host)
		}
	}
	q = q.Arg("ports", ports)

	return &Service{
		Query:  q,
		Client: r.Client,
	}
}

// Sets a secret given a user-defined name and the file path on the host, and returns the secret.
//
// The file is limited to a size of 512000 bytes.
func (r *Host) SetSecretFile(name string, path string) *Secret {
	q := r.Query.Select("setSecretFile")
	q = q.Arg("name", name)
	q = q.Arg("path", path)

	return &Secret{
		Query:  q,
		Client: r.Client,
	}
}

// HostTunnelOpts contains options for Host.Tunnel
type HostTunnelOpts struct {
	// Configure explicit port forwarding rules for the tunnel.
	//
	// If a port's frontend is unspecified or 0, a random port will be chosen by the host.
	//
	// If no ports are given, all of the service's ports are forwarded. If native is true, each port maps to the same port on the host. If native is false, each port maps to a random port chosen by the host.
	//
	// If ports are given and native is true, the ports are additive.
	Ports []PortForward
	// Map each service port to the same port on the host, as if the service were running natively.
	//
	// Note: enabling may result in port conflicts.
	Native bool
}

// Creates a tunnel that forwards traffic from the host to a service.
func (r *Host) Tunnel(service *Service, opts ...HostTunnelOpts) *Service {
	assertNotNil("service", service)
	q := r.Query.Select("tunnel")
	for i := len(opts) - 1; i >= 0; i-- {
		// `ports` optional argument
		if !querybuilder.IsZeroValue(opts[i].Ports) {
			q = q.Arg("ports", opts[i].Ports)
		}
		// `native` optional argument
		if !querybuilder.IsZeroValue(opts[i].Native) {
			q = q.Arg("native", opts[i].Native)
		}
	}
	q = q.Arg("service", service)

	return &Service{
		Query:  q,
		Client: r.Client,
	}
}

// Accesses a Unix socket on the host.
func (r *Host) UnixSocket(path string) *Socket {
	q := r.Query.Select("unixSocket")
	q = q.Arg("path", path)

	return &Socket{
		Query:  q,
		Client: r.Client,
	}
}

// A graphql input type, which is essentially just a group of named args.
// This is currently only used to represent pre-existing usage of graphql input types
// in the core API. It is not used by user modules and shouldn't ever be as user
// module accept input objects via their id rather than graphql input types.
type InputTypeDef struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	id   *InputTypeDefID
	name *string
}

// Static fields defined on this input object, if any.
func (r *InputTypeDef) Fields(ctx context.Context) ([]FieldTypeDef, error) {
	q := r.Query.Select("fields")

	q = q.Select("id")

	type fields struct {
		Id FieldTypeDefID
	}

	convert := func(fields []fields) []FieldTypeDef {
		out := []FieldTypeDef{}

		for i := range fields {
			val := FieldTypeDef{id: &fields[i].Id}
			val.Query = querybuilder.Query().Select("loadFieldTypeDefFromID").Arg("id", fields[i].Id)
			val.Client = r.Client
			out = append(out, val)
		}

		return out
	}
	var response []fields

	q = q.Bind(&response)

	err := q.Execute(ctx, r.Client)
	if err != nil {
		return nil, err
	}

	return convert(response), nil
}

// A unique identifier for this InputTypeDef.
func (r *InputTypeDef) ID(ctx context.Context) (InputTypeDefID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response InputTypeDefID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *InputTypeDef) XXX_GraphQLType() string {
	return "InputTypeDef"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *InputTypeDef) XXX_GraphQLIDType() string {
	return "InputTypeDefID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *InputTypeDef) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *InputTypeDef) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The name of the input object.
func (r *InputTypeDef) Name(ctx context.Context) (string, error) {
	if r.name != nil {
		return *r.name, nil
	}
	q := r.Query.Select("name")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A definition of a custom interface defined in a Module.
type InterfaceTypeDef struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	description      *string
	id               *InterfaceTypeDefID
	name             *string
	sourceModuleName *string
}

// The doc string for the interface, if any.
func (r *InterfaceTypeDef) Description(ctx context.Context) (string, error) {
	if r.description != nil {
		return *r.description, nil
	}
	q := r.Query.Select("description")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Functions defined on this interface, if any.
func (r *InterfaceTypeDef) Functions(ctx context.Context) ([]Function, error) {
	q := r.Query.Select("functions")

	q = q.Select("id")

	type functions struct {
		Id FunctionID
	}

	convert := func(fields []functions) []Function {
		out := []Function{}

		for i := range fields {
			val := Function{id: &fields[i].Id}
			val.Query = querybuilder.Query().Select("loadFunctionFromID").Arg("id", fields[i].Id)
			val.Client = r.Client
			out = append(out, val)
		}

		return out
	}
	var response []functions

	q = q.Bind(&response)

	err := q.Execute(ctx, r.Client)
	if err != nil {
		return nil, err
	}

	return convert(response), nil
}

// A unique identifier for this InterfaceTypeDef.
func (r *InterfaceTypeDef) ID(ctx context.Context) (InterfaceTypeDefID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response InterfaceTypeDefID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *InterfaceTypeDef) XXX_GraphQLType() string {
	return "InterfaceTypeDef"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *InterfaceTypeDef) XXX_GraphQLIDType() string {
	return "InterfaceTypeDefID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *InterfaceTypeDef) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *InterfaceTypeDef) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The name of the interface.
func (r *InterfaceTypeDef) Name(ctx context.Context) (string, error) {
	if r.name != nil {
		return *r.name, nil
	}
	q := r.Query.Select("name")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// If this InterfaceTypeDef is associated with a Module, the name of the module. Unset otherwise.
func (r *InterfaceTypeDef) SourceModuleName(ctx context.Context) (string, error) {
	if r.sourceModuleName != nil {
		return *r.sourceModuleName, nil
	}
	q := r.Query.Select("sourceModuleName")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A simple key value object that represents a label.
type Label struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	id    *LabelID
	name  *string
	value *string
}

// A unique identifier for this Label.
func (r *Label) ID(ctx context.Context) (LabelID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response LabelID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *Label) XXX_GraphQLType() string {
	return "Label"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *Label) XXX_GraphQLIDType() string {
	return "LabelID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *Label) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *Label) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The label name.
func (r *Label) Name(ctx context.Context) (string, error) {
	if r.name != nil {
		return *r.name, nil
	}
	q := r.Query.Select("name")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The label value.
func (r *Label) Value(ctx context.Context) (string, error) {
	if r.value != nil {
		return *r.value, nil
	}
	q := r.Query.Select("value")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A definition of a list type in a Module.
type ListTypeDef struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	id *ListTypeDefID
}

// The type of the elements in the list.
func (r *ListTypeDef) ElementTypeDef() *TypeDef {
	q := r.Query.Select("elementTypeDef")

	return &TypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// A unique identifier for this ListTypeDef.
func (r *ListTypeDef) ID(ctx context.Context) (ListTypeDefID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response ListTypeDefID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *ListTypeDef) XXX_GraphQLType() string {
	return "ListTypeDef"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *ListTypeDef) XXX_GraphQLIDType() string {
	return "ListTypeDefID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *ListTypeDef) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *ListTypeDef) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// Module source that that originates from a path locally relative to an arbitrary directory.
type LocalModuleSource struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	id          *LocalModuleSourceID
	rootSubpath *string
}

// The directory containing everything needed to load load and use the module.
func (r *LocalModuleSource) ContextDirectory() *Directory {
	q := r.Query.Select("contextDirectory")

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// A unique identifier for this LocalModuleSource.
func (r *LocalModuleSource) ID(ctx context.Context) (LocalModuleSourceID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response LocalModuleSourceID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *LocalModuleSource) XXX_GraphQLType() string {
	return "LocalModuleSource"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *LocalModuleSource) XXX_GraphQLIDType() string {
	return "LocalModuleSourceID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *LocalModuleSource) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *LocalModuleSource) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The path to the root of the module source under the context directory. This directory contains its configuration file. It also contains its source code (possibly as a subdirectory).
func (r *LocalModuleSource) RootSubpath(ctx context.Context) (string, error) {
	if r.rootSubpath != nil {
		return *r.rootSubpath, nil
	}
	q := r.Query.Select("rootSubpath")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A Dagger module.
type Module struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	description *string
	id          *ModuleID
	name        *string
	sdk         *string
	serve       *Void
}
type WithModuleFunc func(r *Module) *Module

// With calls the provided function with current Module.
//
// This is useful for reusability and readability by not breaking the calling chain.
func (r *Module) With(f WithModuleFunc) *Module {
	return f(r)
}

// Modules used by this module.
func (r *Module) Dependencies(ctx context.Context) ([]Module, error) {
	q := r.Query.Select("dependencies")

	q = q.Select("id")

	type dependencies struct {
		Id ModuleID
	}

	convert := func(fields []dependencies) []Module {
		out := []Module{}

		for i := range fields {
			val := Module{id: &fields[i].Id}
			val.Query = querybuilder.Query().Select("loadModuleFromID").Arg("id", fields[i].Id)
			val.Client = r.Client
			out = append(out, val)
		}

		return out
	}
	var response []dependencies

	q = q.Bind(&response)

	err := q.Execute(ctx, r.Client)
	if err != nil {
		return nil, err
	}

	return convert(response), nil
}

// The dependencies as configured by the module.
func (r *Module) DependencyConfig(ctx context.Context) ([]ModuleDependency, error) {
	q := r.Query.Select("dependencyConfig")

	q = q.Select("id")

	type dependencyConfig struct {
		Id ModuleDependencyID
	}

	convert := func(fields []dependencyConfig) []ModuleDependency {
		out := []ModuleDependency{}

		for i := range fields {
			val := ModuleDependency{id: &fields[i].Id}
			val.Query = querybuilder.Query().Select("loadModuleDependencyFromID").Arg("id", fields[i].Id)
			val.Client = r.Client
			out = append(out, val)
		}

		return out
	}
	var response []dependencyConfig

	q = q.Bind(&response)

	err := q.Execute(ctx, r.Client)
	if err != nil {
		return nil, err
	}

	return convert(response), nil
}

// The doc string of the module, if any
func (r *Module) Description(ctx context.Context) (string, error) {
	if r.description != nil {
		return *r.description, nil
	}
	q := r.Query.Select("description")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The generated files and directories made on top of the module source's context directory.
func (r *Module) GeneratedContextDiff() *Directory {
	q := r.Query.Select("generatedContextDiff")

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// The module source's context plus any configuration and source files created by codegen.
func (r *Module) GeneratedContextDirectory() *Directory {
	q := r.Query.Select("generatedContextDirectory")

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// A unique identifier for this Module.
func (r *Module) ID(ctx context.Context) (ModuleID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response ModuleID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *Module) XXX_GraphQLType() string {
	return "Module"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *Module) XXX_GraphQLIDType() string {
	return "ModuleID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *Module) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *Module) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// Retrieves the module with the objects loaded via its SDK.
func (r *Module) Initialize() *Module {
	q := r.Query.Select("initialize")

	return &Module{
		Query:  q,
		Client: r.Client,
	}
}

// Interfaces served by this module.
func (r *Module) Interfaces(ctx context.Context) ([]TypeDef, error) {
	q := r.Query.Select("interfaces")

	q = q.Select("id")

	type interfaces struct {
		Id TypeDefID
	}

	convert := func(fields []interfaces) []TypeDef {
		out := []TypeDef{}

		for i := range fields {
			val := TypeDef{id: &fields[i].Id}
			val.Query = querybuilder.Query().Select("loadTypeDefFromID").Arg("id", fields[i].Id)
			val.Client = r.Client
			out = append(out, val)
		}

		return out
	}
	var response []interfaces

	q = q.Bind(&response)

	err := q.Execute(ctx, r.Client)
	if err != nil {
		return nil, err
	}

	return convert(response), nil
}

// The name of the module
func (r *Module) Name(ctx context.Context) (string, error) {
	if r.name != nil {
		return *r.name, nil
	}
	q := r.Query.Select("name")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Objects served by this module.
func (r *Module) Objects(ctx context.Context) ([]TypeDef, error) {
	q := r.Query.Select("objects")

	q = q.Select("id")

	type objects struct {
		Id TypeDefID
	}

	convert := func(fields []objects) []TypeDef {
		out := []TypeDef{}

		for i := range fields {
			val := TypeDef{id: &fields[i].Id}
			val.Query = querybuilder.Query().Select("loadTypeDefFromID").Arg("id", fields[i].Id)
			val.Client = r.Client
			out = append(out, val)
		}

		return out
	}
	var response []objects

	q = q.Bind(&response)

	err := q.Execute(ctx, r.Client)
	if err != nil {
		return nil, err
	}

	return convert(response), nil
}

// The container that runs the module's entrypoint. It will fail to execute if the module doesn't compile.
func (r *Module) Runtime() *Container {
	q := r.Query.Select("runtime")

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// The SDK used by this module. Either a name of a builtin SDK or a module source ref string pointing to the SDK's implementation.
func (r *Module) SDK(ctx context.Context) (string, error) {
	if r.sdk != nil {
		return *r.sdk, nil
	}
	q := r.Query.Select("sdk")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Serve a module's API in the current session.
//
// Note: this can only be called once per session. In the future, it could return a stream or service to remove the side effect.
func (r *Module) Serve(ctx context.Context) (Void, error) {
	progParent := progrock.FromContext(ctx).Parent
	progrock.FromContext(ctx).Warn("Serve propagating parent", progrock.Labelf("parent", progParent))
	if r.serve != nil {
		return *r.serve, nil
	}
	q := r.Query.Select("serve")

	var response Void

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The source for the module.
func (r *Module) Source() *ModuleSource {
	q := r.Query.Select("source")

	return &ModuleSource{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves the module with the given description
func (r *Module) WithDescription(description string) *Module {
	q := r.Query.Select("withDescription")
	q = q.Arg("description", description)

	return &Module{
		Query:  q,
		Client: r.Client,
	}
}

// This module plus the given Interface type and associated functions
func (r *Module) WithInterface(iface *TypeDef) *Module {
	assertNotNil("iface", iface)
	q := r.Query.Select("withInterface")
	q = q.Arg("iface", iface)

	return &Module{
		Query:  q,
		Client: r.Client,
	}
}

// This module plus the given Object type and associated functions.
func (r *Module) WithObject(object *TypeDef) *Module {
	assertNotNil("object", object)
	q := r.Query.Select("withObject")
	q = q.Arg("object", object)

	return &Module{
		Query:  q,
		Client: r.Client,
	}
}

// Retrieves the module with basic configuration loaded if present.
func (r *Module) WithSource(source *ModuleSource) *Module {
	assertNotNil("source", source)
	q := r.Query.Select("withSource")
	q = q.Arg("source", source)

	return &Module{
		Query:  q,
		Client: r.Client,
	}
}

// The configuration of dependency of a module.
type ModuleDependency struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	id   *ModuleDependencyID
	name *string
}

// A unique identifier for this ModuleDependency.
func (r *ModuleDependency) ID(ctx context.Context) (ModuleDependencyID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response ModuleDependencyID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *ModuleDependency) XXX_GraphQLType() string {
	return "ModuleDependency"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *ModuleDependency) XXX_GraphQLIDType() string {
	return "ModuleDependencyID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *ModuleDependency) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *ModuleDependency) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The name of the dependency module.
func (r *ModuleDependency) Name(ctx context.Context) (string, error) {
	if r.name != nil {
		return *r.name, nil
	}
	q := r.Query.Select("name")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The source for the dependency module.
func (r *ModuleDependency) Source() *ModuleSource {
	q := r.Query.Select("source")

	return &ModuleSource{
		Query:  q,
		Client: r.Client,
	}
}

// The source needed to load and run a module, along with any metadata about the source such as versions/urls/etc.
type ModuleSource struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	asString                     *string
	configExists                 *bool
	id                           *ModuleSourceID
	kind                         *ModuleSourceKind
	moduleName                   *string
	moduleOriginalName           *string
	resolveContextPathFromCaller *string
	sourceRootSubpath            *string
	sourceSubpath                *string
}
type WithModuleSourceFunc func(r *ModuleSource) *ModuleSource

// With calls the provided function with current ModuleSource.
//
// This is useful for reusability and readability by not breaking the calling chain.
func (r *ModuleSource) With(f WithModuleSourceFunc) *ModuleSource {
	return f(r)
}

// If the source is a of kind git, the git source representation of it.
func (r *ModuleSource) AsGitSource() *GitModuleSource {
	q := r.Query.Select("asGitSource")

	return &GitModuleSource{
		Query:  q,
		Client: r.Client,
	}
}

// If the source is of kind local, the local source representation of it.
func (r *ModuleSource) AsLocalSource() *LocalModuleSource {
	q := r.Query.Select("asLocalSource")

	return &LocalModuleSource{
		Query:  q,
		Client: r.Client,
	}
}

// Load the source as a module. If this is a local source, the parent directory must have been provided during module source creation
func (r *ModuleSource) AsModule() *Module {
	q := r.Query.Select("asModule")

	return &Module{
		Query:  q,
		Client: r.Client,
	}
}

// A human readable ref string representation of this module source.
func (r *ModuleSource) AsString(ctx context.Context) (string, error) {
	if r.asString != nil {
		return *r.asString, nil
	}
	q := r.Query.Select("asString")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Returns whether the module source has a configuration file.
func (r *ModuleSource) ConfigExists(ctx context.Context) (bool, error) {
	if r.configExists != nil {
		return *r.configExists, nil
	}
	q := r.Query.Select("configExists")

	var response bool

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The directory containing everything needed to load load and use the module.
func (r *ModuleSource) ContextDirectory() *Directory {
	q := r.Query.Select("contextDirectory")

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// The dependencies of the module source. Includes dependencies from the configuration and any extras from withDependencies calls.
func (r *ModuleSource) Dependencies(ctx context.Context) ([]ModuleDependency, error) {
	q := r.Query.Select("dependencies")

	q = q.Select("id")

	type dependencies struct {
		Id ModuleDependencyID
	}

	convert := func(fields []dependencies) []ModuleDependency {
		out := []ModuleDependency{}

		for i := range fields {
			val := ModuleDependency{id: &fields[i].Id}
			val.Query = querybuilder.Query().Select("loadModuleDependencyFromID").Arg("id", fields[i].Id)
			val.Client = r.Client
			out = append(out, val)
		}

		return out
	}
	var response []dependencies

	q = q.Bind(&response)

	err := q.Execute(ctx, r.Client)
	if err != nil {
		return nil, err
	}

	return convert(response), nil
}

// The directory containing the module configuration and source code (source code may be in a subdir).
func (r *ModuleSource) Directory(path string) *Directory {
	q := r.Query.Select("directory")
	q = q.Arg("path", path)

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// A unique identifier for this ModuleSource.
func (r *ModuleSource) ID(ctx context.Context) (ModuleSourceID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response ModuleSourceID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *ModuleSource) XXX_GraphQLType() string {
	return "ModuleSource"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *ModuleSource) XXX_GraphQLIDType() string {
	return "ModuleSourceID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *ModuleSource) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *ModuleSource) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The kind of source (e.g. local, git, etc.)
func (r *ModuleSource) Kind(ctx context.Context) (ModuleSourceKind, error) {
	if r.kind != nil {
		return *r.kind, nil
	}
	q := r.Query.Select("kind")

	var response ModuleSourceKind

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// If set, the name of the module this source references, including any overrides at runtime by callers.
func (r *ModuleSource) ModuleName(ctx context.Context) (string, error) {
	if r.moduleName != nil {
		return *r.moduleName, nil
	}
	q := r.Query.Select("moduleName")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The original name of the module this source references, as defined in the module configuration.
func (r *ModuleSource) ModuleOriginalName(ctx context.Context) (string, error) {
	if r.moduleOriginalName != nil {
		return *r.moduleOriginalName, nil
	}
	q := r.Query.Select("moduleOriginalName")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The path to the module source's context directory on the caller's filesystem. Only valid for local sources.
func (r *ModuleSource) ResolveContextPathFromCaller(ctx context.Context) (string, error) {
	if r.resolveContextPathFromCaller != nil {
		return *r.resolveContextPathFromCaller, nil
	}
	q := r.Query.Select("resolveContextPathFromCaller")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Resolve the provided module source arg as a dependency relative to this module source.
func (r *ModuleSource) ResolveDependency(dep *ModuleSource) *ModuleSource {
	assertNotNil("dep", dep)
	q := r.Query.Select("resolveDependency")
	q = q.Arg("dep", dep)

	return &ModuleSource{
		Query:  q,
		Client: r.Client,
	}
}

// Load the source from its path on the caller's filesystem, including only needed+configured files and directories. Only valid for local sources.
func (r *ModuleSource) ResolveFromCaller() *ModuleSource {
	q := r.Query.Select("resolveFromCaller")

	return &ModuleSource{
		Query:  q,
		Client: r.Client,
	}
}

// The path relative to context of the root of the module source, which contains dagger.json. It also contains the module implementation source code, but that may or may not being a subdir of this root.
func (r *ModuleSource) SourceRootSubpath(ctx context.Context) (string, error) {
	if r.sourceRootSubpath != nil {
		return *r.sourceRootSubpath, nil
	}
	q := r.Query.Select("sourceRootSubpath")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The path relative to context of the module implementation source code.
func (r *ModuleSource) SourceSubpath(ctx context.Context) (string, error) {
	if r.sourceSubpath != nil {
		return *r.sourceSubpath, nil
	}
	q := r.Query.Select("sourceSubpath")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Update the module source with a new context directory. Only valid for local sources.
func (r *ModuleSource) WithContextDirectory(dir *Directory) *ModuleSource {
	assertNotNil("dir", dir)
	q := r.Query.Select("withContextDirectory")
	q = q.Arg("dir", dir)

	return &ModuleSource{
		Query:  q,
		Client: r.Client,
	}
}

// Append the provided dependencies to the module source's dependency list.
func (r *ModuleSource) WithDependencies(dependencies []*ModuleDependency) *ModuleSource {
	q := r.Query.Select("withDependencies")
	q = q.Arg("dependencies", dependencies)

	return &ModuleSource{
		Query:  q,
		Client: r.Client,
	}
}

// Update the module source with a new name.
func (r *ModuleSource) WithName(name string) *ModuleSource {
	q := r.Query.Select("withName")
	q = q.Arg("name", name)

	return &ModuleSource{
		Query:  q,
		Client: r.Client,
	}
}

// Update the module source with a new SDK.
func (r *ModuleSource) WithSDK(sdk string) *ModuleSource {
	q := r.Query.Select("withSDK")
	q = q.Arg("sdk", sdk)

	return &ModuleSource{
		Query:  q,
		Client: r.Client,
	}
}

// Update the module source with a new source subpath.
func (r *ModuleSource) WithSourceSubpath(path string) *ModuleSource {
	q := r.Query.Select("withSourceSubpath")
	q = q.Arg("path", path)

	return &ModuleSource{
		Query:  q,
		Client: r.Client,
	}
}

// A definition of a custom object defined in a Module.
type ObjectTypeDef struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	description      *string
	id               *ObjectTypeDefID
	name             *string
	sourceModuleName *string
}

// The function used to construct new instances of this object, if any
func (r *ObjectTypeDef) Constructor() *Function {
	q := r.Query.Select("constructor")

	return &Function{
		Query:  q,
		Client: r.Client,
	}
}

// The doc string for the object, if any.
func (r *ObjectTypeDef) Description(ctx context.Context) (string, error) {
	if r.description != nil {
		return *r.description, nil
	}
	q := r.Query.Select("description")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Static fields defined on this object, if any.
func (r *ObjectTypeDef) Fields(ctx context.Context) ([]FieldTypeDef, error) {
	q := r.Query.Select("fields")

	q = q.Select("id")

	type fields struct {
		Id FieldTypeDefID
	}

	convert := func(fields []fields) []FieldTypeDef {
		out := []FieldTypeDef{}

		for i := range fields {
			val := FieldTypeDef{id: &fields[i].Id}
			val.Query = querybuilder.Query().Select("loadFieldTypeDefFromID").Arg("id", fields[i].Id)
			val.Client = r.Client
			out = append(out, val)
		}

		return out
	}
	var response []fields

	q = q.Bind(&response)

	err := q.Execute(ctx, r.Client)
	if err != nil {
		return nil, err
	}

	return convert(response), nil
}

// Functions defined on this object, if any.
func (r *ObjectTypeDef) Functions(ctx context.Context) ([]Function, error) {
	q := r.Query.Select("functions")

	q = q.Select("id")

	type functions struct {
		Id FunctionID
	}

	convert := func(fields []functions) []Function {
		out := []Function{}

		for i := range fields {
			val := Function{id: &fields[i].Id}
			val.Query = querybuilder.Query().Select("loadFunctionFromID").Arg("id", fields[i].Id)
			val.Client = r.Client
			out = append(out, val)
		}

		return out
	}
	var response []functions

	q = q.Bind(&response)

	err := q.Execute(ctx, r.Client)
	if err != nil {
		return nil, err
	}

	return convert(response), nil
}

// A unique identifier for this ObjectTypeDef.
func (r *ObjectTypeDef) ID(ctx context.Context) (ObjectTypeDefID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response ObjectTypeDefID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *ObjectTypeDef) XXX_GraphQLType() string {
	return "ObjectTypeDef"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *ObjectTypeDef) XXX_GraphQLIDType() string {
	return "ObjectTypeDefID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *ObjectTypeDef) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *ObjectTypeDef) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The name of the object.
func (r *ObjectTypeDef) Name(ctx context.Context) (string, error) {
	if r.name != nil {
		return *r.name, nil
	}
	q := r.Query.Select("name")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// If this ObjectTypeDef is associated with a Module, the name of the module. Unset otherwise.
func (r *ObjectTypeDef) SourceModuleName(ctx context.Context) (string, error) {
	if r.sourceModuleName != nil {
		return *r.sourceModuleName, nil
	}
	q := r.Query.Select("sourceModuleName")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A port exposed by a container.
type Port struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	description                 *string
	experimentalSkipHealthcheck *bool
	id                          *PortID
	port                        *int
	protocol                    *NetworkProtocol
}

// The port description.
func (r *Port) Description(ctx context.Context) (string, error) {
	if r.description != nil {
		return *r.description, nil
	}
	q := r.Query.Select("description")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Skip the health check when run as a service.
func (r *Port) ExperimentalSkipHealthcheck(ctx context.Context) (bool, error) {
	if r.experimentalSkipHealthcheck != nil {
		return *r.experimentalSkipHealthcheck, nil
	}
	q := r.Query.Select("experimentalSkipHealthcheck")

	var response bool

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A unique identifier for this Port.
func (r *Port) ID(ctx context.Context) (PortID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response PortID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *Port) XXX_GraphQLType() string {
	return "Port"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *Port) XXX_GraphQLIDType() string {
	return "PortID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *Port) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *Port) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The port number.
func (r *Port) Port(ctx context.Context) (int, error) {
	if r.port != nil {
		return *r.port, nil
	}
	q := r.Query.Select("port")

	var response int

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// The transport layer protocol.
func (r *Port) Protocol(ctx context.Context) (NetworkProtocol, error) {
	if r.protocol != nil {
		return *r.protocol, nil
	}
	q := r.Query.Select("protocol")

	var response NetworkProtocol

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

type WithClientFunc func(r *Client) *Client

// With calls the provided function with current Client.
//
// This is useful for reusability and readability by not breaking the calling chain.
func (r *Client) With(f WithClientFunc) *Client {
	return f(r)
}

// Retrieves a content-addressed blob.
func (r *Client) Blob(digest string, size int, mediaType string, uncompressed string) *Directory {
	q := r.Query.Select("blob")
	q = q.Arg("digest", digest)
	q = q.Arg("size", size)
	q = q.Arg("mediaType", mediaType)
	q = q.Arg("uncompressed", uncompressed)

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// Constructs a cache volume for a given cache key.
func (r *Client) CacheVolume(key string) *CacheVolume {
	q := r.Query.Select("cacheVolume")
	q = q.Arg("key", key)

	return &CacheVolume{
		Query:  q,
		Client: r.Client,
	}
}

// Checks if the current Dagger Engine is compatible with an SDK's required version.
func (r *Client) CheckVersionCompatibility(ctx context.Context, version string) (bool, error) {
	q := r.Query.Select("checkVersionCompatibility")
	q = q.Arg("version", version)

	var response bool

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// ContainerOpts contains options for Client.Container
type ContainerOpts struct {
	// DEPRECATED: Use `loadContainerFromID` instead.
	ID ContainerID
	// Platform to initialize the container with.
	Platform Platform
}

// Creates a scratch container.
//
// Optional platform argument initializes new containers to execute and publish as that platform. Platform defaults to that of the builder's host.
func (r *Client) Container(opts ...ContainerOpts) *Container {
	q := r.Query.Select("container")
	for i := len(opts) - 1; i >= 0; i-- {
		// `id` optional argument
		if !querybuilder.IsZeroValue(opts[i].ID) {
			q = q.Arg("id", opts[i].ID)
		}
		// `platform` optional argument
		if !querybuilder.IsZeroValue(opts[i].Platform) {
			q = q.Arg("platform", opts[i].Platform)
		}
	}

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// The FunctionCall context that the SDK caller is currently executing in.
//
// If the caller is not currently executing in a function, this will return an error.
func (r *Client) CurrentFunctionCall() *FunctionCall {
	q := r.Query.Select("currentFunctionCall")

	return &FunctionCall{
		Query:  q,
		Client: r.Client,
	}
}

// The module currently being served in the session, if any.
func (r *Client) CurrentModule() *CurrentModule {
	q := r.Query.Select("currentModule")

	return &CurrentModule{
		Query:  q,
		Client: r.Client,
	}
}

// The TypeDef representations of the objects currently being served in the session.
func (r *Client) CurrentTypeDefs(ctx context.Context) ([]TypeDef, error) {
	q := r.Query.Select("currentTypeDefs")

	q = q.Select("id")

	type currentTypeDefs struct {
		Id TypeDefID
	}

	convert := func(fields []currentTypeDefs) []TypeDef {
		out := []TypeDef{}

		for i := range fields {
			val := TypeDef{id: &fields[i].Id}
			val.Query = querybuilder.Query().Select("loadTypeDefFromID").Arg("id", fields[i].Id)
			val.Client = r.Client
			out = append(out, val)
		}

		return out
	}
	var response []currentTypeDefs

	q = q.Bind(&response)

	err := q.Execute(ctx, r.Client)
	if err != nil {
		return nil, err
	}

	return convert(response), nil
}

// The default platform of the engine.
func (r *Client) DefaultPlatform(ctx context.Context) (Platform, error) {
	q := r.Query.Select("defaultPlatform")

	var response Platform

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// DirectoryOpts contains options for Client.Directory
type DirectoryOpts struct {
	// DEPRECATED: Use `loadDirectoryFromID` isntead.
	ID DirectoryID
}

// Creates an empty directory.
func (r *Client) Directory(opts ...DirectoryOpts) *Directory {
	q := r.Query.Select("directory")
	for i := len(opts) - 1; i >= 0; i-- {
		// `id` optional argument
		if !querybuilder.IsZeroValue(opts[i].ID) {
			q = q.Arg("id", opts[i].ID)
		}
	}

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// Deprecated: Use LoadFileFromID instead.
func (r *Client) File(id FileID) *File {
	q := r.Query.Select("file")
	q = q.Arg("id", id)

	return &File{
		Query:  q,
		Client: r.Client,
	}
}

// Creates a function.
func (r *Client) Function(name string, returnType *TypeDef) *Function {
	assertNotNil("returnType", returnType)
	q := r.Query.Select("function")
	q = q.Arg("name", name)
	q = q.Arg("returnType", returnType)

	return &Function{
		Query:  q,
		Client: r.Client,
	}
}

// Create a code generation result, given a directory containing the generated code.
func (r *Client) GeneratedCode(code *Directory) *GeneratedCode {
	assertNotNil("code", code)
	q := r.Query.Select("generatedCode")
	q = q.Arg("code", code)

	return &GeneratedCode{
		Query:  q,
		Client: r.Client,
	}
}

// GitOpts contains options for Client.Git
type GitOpts struct {
	// Set to true to keep .git directory.
	KeepGitDir bool
	// A service which must be started before the repo is fetched.
	ExperimentalServiceHost *Service
	// Set SSH known hosts
	SSHKnownHosts string
	// Set SSH auth socket
	SSHAuthSocket *Socket
}

// Queries a Git repository.
func (r *Client) Git(url string, opts ...GitOpts) *GitRepository {
	q := r.Query.Select("git")
	for i := len(opts) - 1; i >= 0; i-- {
		// `keepGitDir` optional argument
		if !querybuilder.IsZeroValue(opts[i].KeepGitDir) {
			q = q.Arg("keepGitDir", opts[i].KeepGitDir)
		}
		// `experimentalServiceHost` optional argument
		if !querybuilder.IsZeroValue(opts[i].ExperimentalServiceHost) {
			q = q.Arg("experimentalServiceHost", opts[i].ExperimentalServiceHost)
		}
		// `sshKnownHosts` optional argument
		if !querybuilder.IsZeroValue(opts[i].SSHKnownHosts) {
			q = q.Arg("sshKnownHosts", opts[i].SSHKnownHosts)
		}
		// `sshAuthSocket` optional argument
		if !querybuilder.IsZeroValue(opts[i].SSHAuthSocket) {
			q = q.Arg("sshAuthSocket", opts[i].SSHAuthSocket)
		}
	}
	q = q.Arg("url", url)

	return &GitRepository{
		Query:  q,
		Client: r.Client,
	}
}

// Queries the host environment.
func (r *Client) Host() *Host {
	q := r.Query.Select("host")

	return &Host{
		Query:  q,
		Client: r.Client,
	}
}

// HTTPOpts contains options for Client.HTTP
type HTTPOpts struct {
	// A service which must be started before the URL is fetched.
	ExperimentalServiceHost *Service
}

// Returns a file containing an http remote url content.
func (r *Client) HTTP(url string, opts ...HTTPOpts) *File {
	q := r.Query.Select("http")
	for i := len(opts) - 1; i >= 0; i-- {
		// `experimentalServiceHost` optional argument
		if !querybuilder.IsZeroValue(opts[i].ExperimentalServiceHost) {
			q = q.Arg("experimentalServiceHost", opts[i].ExperimentalServiceHost)
		}
	}
	q = q.Arg("url", url)

	return &File{
		Query:  q,
		Client: r.Client,
	}
}

// Load a CacheVolume from its ID.
func (r *Client) LoadCacheVolumeFromID(id CacheVolumeID) *CacheVolume {
	q := r.Query.Select("loadCacheVolumeFromID")
	q = q.Arg("id", id)

	return &CacheVolume{
		Query:  q,
		Client: r.Client,
	}
}

// Load a Container from its ID.
func (r *Client) LoadContainerFromID(id ContainerID) *Container {
	q := r.Query.Select("loadContainerFromID")
	q = q.Arg("id", id)

	return &Container{
		Query:  q,
		Client: r.Client,
	}
}

// Load a CurrentModule from its ID.
func (r *Client) LoadCurrentModuleFromID(id CurrentModuleID) *CurrentModule {
	q := r.Query.Select("loadCurrentModuleFromID")
	q = q.Arg("id", id)

	return &CurrentModule{
		Query:  q,
		Client: r.Client,
	}
}

// Load a Directory from its ID.
func (r *Client) LoadDirectoryFromID(id DirectoryID) *Directory {
	q := r.Query.Select("loadDirectoryFromID")
	q = q.Arg("id", id)

	return &Directory{
		Query:  q,
		Client: r.Client,
	}
}

// Load a EnvVariable from its ID.
func (r *Client) LoadEnvVariableFromID(id EnvVariableID) *EnvVariable {
	q := r.Query.Select("loadEnvVariableFromID")
	q = q.Arg("id", id)

	return &EnvVariable{
		Query:  q,
		Client: r.Client,
	}
}

// Load a FieldTypeDef from its ID.
func (r *Client) LoadFieldTypeDefFromID(id FieldTypeDefID) *FieldTypeDef {
	q := r.Query.Select("loadFieldTypeDefFromID")
	q = q.Arg("id", id)

	return &FieldTypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// Load a File from its ID.
func (r *Client) LoadFileFromID(id FileID) *File {
	q := r.Query.Select("loadFileFromID")
	q = q.Arg("id", id)

	return &File{
		Query:  q,
		Client: r.Client,
	}
}

// Load a FunctionArg from its ID.
func (r *Client) LoadFunctionArgFromID(id FunctionArgID) *FunctionArg {
	q := r.Query.Select("loadFunctionArgFromID")
	q = q.Arg("id", id)

	return &FunctionArg{
		Query:  q,
		Client: r.Client,
	}
}

// Load a FunctionCallArgValue from its ID.
func (r *Client) LoadFunctionCallArgValueFromID(id FunctionCallArgValueID) *FunctionCallArgValue {
	q := r.Query.Select("loadFunctionCallArgValueFromID")
	q = q.Arg("id", id)

	return &FunctionCallArgValue{
		Query:  q,
		Client: r.Client,
	}
}

// Load a FunctionCall from its ID.
func (r *Client) LoadFunctionCallFromID(id FunctionCallID) *FunctionCall {
	q := r.Query.Select("loadFunctionCallFromID")
	q = q.Arg("id", id)

	return &FunctionCall{
		Query:  q,
		Client: r.Client,
	}
}

// Load a Function from its ID.
func (r *Client) LoadFunctionFromID(id FunctionID) *Function {
	q := r.Query.Select("loadFunctionFromID")
	q = q.Arg("id", id)

	return &Function{
		Query:  q,
		Client: r.Client,
	}
}

// Load a GeneratedCode from its ID.
func (r *Client) LoadGeneratedCodeFromID(id GeneratedCodeID) *GeneratedCode {
	q := r.Query.Select("loadGeneratedCodeFromID")
	q = q.Arg("id", id)

	return &GeneratedCode{
		Query:  q,
		Client: r.Client,
	}
}

// Load a GitModuleSource from its ID.
func (r *Client) LoadGitModuleSourceFromID(id GitModuleSourceID) *GitModuleSource {
	q := r.Query.Select("loadGitModuleSourceFromID")
	q = q.Arg("id", id)

	return &GitModuleSource{
		Query:  q,
		Client: r.Client,
	}
}

// Load a GitRef from its ID.
func (r *Client) LoadGitRefFromID(id GitRefID) *GitRef {
	q := r.Query.Select("loadGitRefFromID")
	q = q.Arg("id", id)

	return &GitRef{
		Query:  q,
		Client: r.Client,
	}
}

// Load a GitRepository from its ID.
func (r *Client) LoadGitRepositoryFromID(id GitRepositoryID) *GitRepository {
	q := r.Query.Select("loadGitRepositoryFromID")
	q = q.Arg("id", id)

	return &GitRepository{
		Query:  q,
		Client: r.Client,
	}
}

// Load a Host from its ID.
func (r *Client) LoadHostFromID(id HostID) *Host {
	q := r.Query.Select("loadHostFromID")
	q = q.Arg("id", id)

	return &Host{
		Query:  q,
		Client: r.Client,
	}
}

// Load a InputTypeDef from its ID.
func (r *Client) LoadInputTypeDefFromID(id InputTypeDefID) *InputTypeDef {
	q := r.Query.Select("loadInputTypeDefFromID")
	q = q.Arg("id", id)

	return &InputTypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// Load a InterfaceTypeDef from its ID.
func (r *Client) LoadInterfaceTypeDefFromID(id InterfaceTypeDefID) *InterfaceTypeDef {
	q := r.Query.Select("loadInterfaceTypeDefFromID")
	q = q.Arg("id", id)

	return &InterfaceTypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// Load a Label from its ID.
func (r *Client) LoadLabelFromID(id LabelID) *Label {
	q := r.Query.Select("loadLabelFromID")
	q = q.Arg("id", id)

	return &Label{
		Query:  q,
		Client: r.Client,
	}
}

// Load a ListTypeDef from its ID.
func (r *Client) LoadListTypeDefFromID(id ListTypeDefID) *ListTypeDef {
	q := r.Query.Select("loadListTypeDefFromID")
	q = q.Arg("id", id)

	return &ListTypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// Load a LocalModuleSource from its ID.
func (r *Client) LoadLocalModuleSourceFromID(id LocalModuleSourceID) *LocalModuleSource {
	q := r.Query.Select("loadLocalModuleSourceFromID")
	q = q.Arg("id", id)

	return &LocalModuleSource{
		Query:  q,
		Client: r.Client,
	}
}

// Load a ModuleDependency from its ID.
func (r *Client) LoadModuleDependencyFromID(id ModuleDependencyID) *ModuleDependency {
	q := r.Query.Select("loadModuleDependencyFromID")
	q = q.Arg("id", id)

	return &ModuleDependency{
		Query:  q,
		Client: r.Client,
	}
}

// Load a Module from its ID.
func (r *Client) LoadModuleFromID(id ModuleID) *Module {
	q := r.Query.Select("loadModuleFromID")
	q = q.Arg("id", id)

	return &Module{
		Query:  q,
		Client: r.Client,
	}
}

// Load a ModuleSource from its ID.
func (r *Client) LoadModuleSourceFromID(id ModuleSourceID) *ModuleSource {
	q := r.Query.Select("loadModuleSourceFromID")
	q = q.Arg("id", id)

	return &ModuleSource{
		Query:  q,
		Client: r.Client,
	}
}

// Load a ObjectTypeDef from its ID.
func (r *Client) LoadObjectTypeDefFromID(id ObjectTypeDefID) *ObjectTypeDef {
	q := r.Query.Select("loadObjectTypeDefFromID")
	q = q.Arg("id", id)

	return &ObjectTypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// Load a Port from its ID.
func (r *Client) LoadPortFromID(id PortID) *Port {
	q := r.Query.Select("loadPortFromID")
	q = q.Arg("id", id)

	return &Port{
		Query:  q,
		Client: r.Client,
	}
}

// Load a Secret from its ID.
func (r *Client) LoadSecretFromID(id SecretID) *Secret {
	q := r.Query.Select("loadSecretFromID")
	q = q.Arg("id", id)

	return &Secret{
		Query:  q,
		Client: r.Client,
	}
}

// Load a Service from its ID.
func (r *Client) LoadServiceFromID(id ServiceID) *Service {
	q := r.Query.Select("loadServiceFromID")
	q = q.Arg("id", id)

	return &Service{
		Query:  q,
		Client: r.Client,
	}
}

// Load a Socket from its ID.
func (r *Client) LoadSocketFromID(id SocketID) *Socket {
	q := r.Query.Select("loadSocketFromID")
	q = q.Arg("id", id)

	return &Socket{
		Query:  q,
		Client: r.Client,
	}
}

// Load a Terminal from its ID.
func (r *Client) LoadTerminalFromID(id TerminalID) *Terminal {
	q := r.Query.Select("loadTerminalFromID")
	q = q.Arg("id", id)

	return &Terminal{
		Query:  q,
		Client: r.Client,
	}
}

// Load a TypeDef from its ID.
func (r *Client) LoadTypeDefFromID(id TypeDefID) *TypeDef {
	q := r.Query.Select("loadTypeDefFromID")
	q = q.Arg("id", id)

	return &TypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// Create a new module.
func (r *Client) Module() *Module {
	q := r.Query.Select("module")

	return &Module{
		Query:  q,
		Client: r.Client,
	}
}

// ModuleDependencyOpts contains options for Client.ModuleDependency
type ModuleDependencyOpts struct {
	// If set, the name to use for the dependency. Otherwise, once installed to a parent module, the name of the dependency module will be used by default.
	Name string
}

// Create a new module dependency configuration from a module source and name
func (r *Client) ModuleDependency(source *ModuleSource, opts ...ModuleDependencyOpts) *ModuleDependency {
	assertNotNil("source", source)
	q := r.Query.Select("moduleDependency")
	for i := len(opts) - 1; i >= 0; i-- {
		// `name` optional argument
		if !querybuilder.IsZeroValue(opts[i].Name) {
			q = q.Arg("name", opts[i].Name)
		}
	}
	q = q.Arg("source", source)

	return &ModuleDependency{
		Query:  q,
		Client: r.Client,
	}
}

// ModuleSourceOpts contains options for Client.ModuleSource
type ModuleSourceOpts struct {
	// If true, enforce that the source is a stable version for source kinds that support versioning.
	Stable bool
}

// Create a new module source instance from a source ref string.
func (r *Client) ModuleSource(refString string, opts ...ModuleSourceOpts) *ModuleSource {
	q := r.Query.Select("moduleSource")
	for i := len(opts) - 1; i >= 0; i-- {
		// `stable` optional argument
		if !querybuilder.IsZeroValue(opts[i].Stable) {
			q = q.Arg("stable", opts[i].Stable)
		}
	}
	q = q.Arg("refString", refString)

	return &ModuleSource{
		Query:  q,
		Client: r.Client,
	}
}

// PipelineOpts contains options for Client.Pipeline
type PipelineOpts struct {
	// Description of the sub-pipeline.
	Description string
	// Labels to apply to the sub-pipeline.
	Labels []PipelineLabel
}

// Creates a named sub-pipeline.
func (r *Client) Pipeline(name string, opts ...PipelineOpts) *Client {
	q := r.Query.Select("pipeline")
	for i := len(opts) - 1; i >= 0; i-- {
		// `description` optional argument
		if !querybuilder.IsZeroValue(opts[i].Description) {
			q = q.Arg("description", opts[i].Description)
		}
		// `labels` optional argument
		if !querybuilder.IsZeroValue(opts[i].Labels) {
			q = q.Arg("labels", opts[i].Labels)
		}
	}
	q = q.Arg("name", name)

	return &Client{
		Query:  q,
		Client: r.Client,
	}
}

// Reference a secret by name.
func (r *Client) Secret(name string) *Secret {
	q := r.Query.Select("secret")
	q = q.Arg("name", name)

	return &Secret{
		Query:  q,
		Client: r.Client,
	}
}

// Sets a secret given a user defined name to its plaintext and returns the secret.
//
// The plaintext value is limited to a size of 128000 bytes.
func (r *Client) SetSecret(name string, plaintext string) *Secret {
	q := r.Query.Select("setSecret")
	q = q.Arg("name", name)
	q = q.Arg("plaintext", plaintext)

	return &Secret{
		Query:  q,
		Client: r.Client,
	}
}

// Loads a socket by its ID.
//
// Deprecated: Use LoadSocketFromID instead.
func (r *Client) Socket(id SocketID) *Socket {
	q := r.Query.Select("socket")
	q = q.Arg("id", id)

	return &Socket{
		Query:  q,
		Client: r.Client,
	}
}

// Create a new TypeDef.
func (r *Client) TypeDef() *TypeDef {
	q := r.Query.Select("typeDef")

	return &TypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// A reference to a secret value, which can be handled more safely than the value itself.
type Secret struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	id        *SecretID
	plaintext *string
}

// A unique identifier for this Secret.
func (r *Secret) ID(ctx context.Context) (SecretID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response SecretID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *Secret) XXX_GraphQLType() string {
	return "Secret"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *Secret) XXX_GraphQLIDType() string {
	return "SecretID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *Secret) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *Secret) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The value of this secret.
func (r *Secret) Plaintext(ctx context.Context) (string, error) {
	if r.plaintext != nil {
		return *r.plaintext, nil
	}
	q := r.Query.Select("plaintext")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A content-addressed service providing TCP connectivity.
type Service struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	endpoint *string
	hostname *string
	id       *ServiceID
	start    *ServiceID
	stop     *ServiceID
	up       *Void
}

// ServiceEndpointOpts contains options for Service.Endpoint
type ServiceEndpointOpts struct {
	// The exposed port number for the endpoint
	Port int
	// Return a URL with the given scheme, eg. http for http://
	Scheme string
}

// Retrieves an endpoint that clients can use to reach this container.
//
// If no port is specified, the first exposed port is used. If none exist an error is returned.
//
// If a scheme is specified, a URL is returned. Otherwise, a host:port pair is returned.
func (r *Service) Endpoint(ctx context.Context, opts ...ServiceEndpointOpts) (string, error) {
	if r.endpoint != nil {
		return *r.endpoint, nil
	}
	q := r.Query.Select("endpoint")
	for i := len(opts) - 1; i >= 0; i-- {
		// `port` optional argument
		if !querybuilder.IsZeroValue(opts[i].Port) {
			q = q.Arg("port", opts[i].Port)
		}
		// `scheme` optional argument
		if !querybuilder.IsZeroValue(opts[i].Scheme) {
			q = q.Arg("scheme", opts[i].Scheme)
		}
	}

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Retrieves a hostname which can be used by clients to reach this container.
func (r *Service) Hostname(ctx context.Context) (string, error) {
	if r.hostname != nil {
		return *r.hostname, nil
	}
	q := r.Query.Select("hostname")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A unique identifier for this Service.
func (r *Service) ID(ctx context.Context) (ServiceID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response ServiceID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *Service) XXX_GraphQLType() string {
	return "Service"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *Service) XXX_GraphQLIDType() string {
	return "ServiceID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *Service) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *Service) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// Retrieves the list of ports provided by the service.
func (r *Service) Ports(ctx context.Context) ([]Port, error) {
	q := r.Query.Select("ports")

	q = q.Select("id")

	type ports struct {
		Id PortID
	}

	convert := func(fields []ports) []Port {
		out := []Port{}

		for i := range fields {
			val := Port{id: &fields[i].Id}
			val.Query = querybuilder.Query().Select("loadPortFromID").Arg("id", fields[i].Id)
			val.Client = r.Client
			out = append(out, val)
		}

		return out
	}
	var response []ports

	q = q.Bind(&response)

	err := q.Execute(ctx, r.Client)
	if err != nil {
		return nil, err
	}

	return convert(response), nil
}

// Start the service and wait for its health checks to succeed.
//
// Services bound to a Container do not need to be manually started.
func (r *Service) Start(ctx context.Context) (*Service, error) {
	q := r.Query.Select("start")

	return r, q.Execute(ctx, r.Client)
}

// ServiceStopOpts contains options for Service.Stop
type ServiceStopOpts struct {
	// Immediately kill the service without waiting for a graceful exit
	Kill bool
}

// Stop the service.
func (r *Service) Stop(ctx context.Context, opts ...ServiceStopOpts) (*Service, error) {
	q := r.Query.Select("stop")
	for i := len(opts) - 1; i >= 0; i-- {
		// `kill` optional argument
		if !querybuilder.IsZeroValue(opts[i].Kill) {
			q = q.Arg("kill", opts[i].Kill)
		}
	}

	return r, q.Execute(ctx, r.Client)
}

// ServiceUpOpts contains options for Service.Up
type ServiceUpOpts struct {
	// List of frontend/backend port mappings to forward.
	//
	// Frontend is the port accepting traffic on the host, backend is the service port.
	Ports []PortForward
	// Bind each tunnel port to a random port on the host.
	Random bool
}

// Creates a tunnel that forwards traffic from the caller's network to this service.
func (r *Service) Up(ctx context.Context, opts ...ServiceUpOpts) (Void, error) {
	if r.up != nil {
		return *r.up, nil
	}
	q := r.Query.Select("up")
	for i := len(opts) - 1; i >= 0; i-- {
		// `ports` optional argument
		if !querybuilder.IsZeroValue(opts[i].Ports) {
			q = q.Arg("ports", opts[i].Ports)
		}
		// `random` optional argument
		if !querybuilder.IsZeroValue(opts[i].Random) {
			q = q.Arg("random", opts[i].Random)
		}
	}

	var response Void

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A Unix or TCP/IP socket that can be mounted into a container.
type Socket struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	id *SocketID
}

// A unique identifier for this Socket.
func (r *Socket) ID(ctx context.Context) (SocketID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response SocketID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *Socket) XXX_GraphQLType() string {
	return "Socket"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *Socket) XXX_GraphQLIDType() string {
	return "SocketID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *Socket) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *Socket) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// An interactive terminal that clients can connect to.
type Terminal struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	id                *TerminalID
	websocketEndpoint *string
}

// A unique identifier for this Terminal.
func (r *Terminal) ID(ctx context.Context) (TerminalID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response TerminalID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *Terminal) XXX_GraphQLType() string {
	return "Terminal"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *Terminal) XXX_GraphQLIDType() string {
	return "TerminalID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *Terminal) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *Terminal) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// An http endpoint at which this terminal can be connected to over a websocket.
func (r *Terminal) WebsocketEndpoint(ctx context.Context) (string, error) {
	if r.websocketEndpoint != nil {
		return *r.websocketEndpoint, nil
	}
	q := r.Query.Select("websocketEndpoint")

	var response string

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// A definition of a parameter or return type in a Module.
type TypeDef struct {
	Query  *querybuilder.Selection
	Client graphql.Client

	id       *TypeDefID
	kind     *TypeDefKind
	optional *bool
}
type WithTypeDefFunc func(r *TypeDef) *TypeDef

// With calls the provided function with current TypeDef.
//
// This is useful for reusability and readability by not breaking the calling chain.
func (r *TypeDef) With(f WithTypeDefFunc) *TypeDef {
	return f(r)
}

// If kind is INPUT, the input-specific type definition. If kind is not INPUT, this will be null.
func (r *TypeDef) AsInput() *InputTypeDef {
	q := r.Query.Select("asInput")

	return &InputTypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// If kind is INTERFACE, the interface-specific type definition. If kind is not INTERFACE, this will be null.
func (r *TypeDef) AsInterface() *InterfaceTypeDef {
	q := r.Query.Select("asInterface")

	return &InterfaceTypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// If kind is LIST, the list-specific type definition. If kind is not LIST, this will be null.
func (r *TypeDef) AsList() *ListTypeDef {
	q := r.Query.Select("asList")

	return &ListTypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// If kind is OBJECT, the object-specific type definition. If kind is not OBJECT, this will be null.
func (r *TypeDef) AsObject() *ObjectTypeDef {
	q := r.Query.Select("asObject")

	return &ObjectTypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// A unique identifier for this TypeDef.
func (r *TypeDef) ID(ctx context.Context) (TypeDefID, error) {
	if r.id != nil {
		return *r.id, nil
	}
	q := r.Query.Select("id")

	var response TypeDefID

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// XXX_GraphQLType is an internal function. It returns the native GraphQL type name
func (r *TypeDef) XXX_GraphQLType() string {
	return "TypeDef"
}

// XXX_GraphQLIDType is an internal function. It returns the native GraphQL type name for the ID of this object
func (r *TypeDef) XXX_GraphQLIDType() string {
	return "TypeDefID"
}

// XXX_GraphQLID is an internal function. It returns the underlying type ID
func (r *TypeDef) XXX_GraphQLID(ctx context.Context) (string, error) {
	id, err := r.ID(ctx)
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (r *TypeDef) MarshalJSON() ([]byte, error) {
	id, err := r.ID(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(id)
}

// The kind of type this is (e.g. primitive, list, object).
func (r *TypeDef) Kind(ctx context.Context) (TypeDefKind, error) {
	if r.kind != nil {
		return *r.kind, nil
	}
	q := r.Query.Select("kind")

	var response TypeDefKind

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Whether this type can be set to null. Defaults to false.
func (r *TypeDef) Optional(ctx context.Context) (bool, error) {
	if r.optional != nil {
		return *r.optional, nil
	}
	q := r.Query.Select("optional")

	var response bool

	q = q.Bind(&response)
	return response, q.Execute(ctx, r.Client)
}

// Adds a function for constructing a new instance of an Object TypeDef, failing if the type is not an object.
func (r *TypeDef) WithConstructor(function *Function) *TypeDef {
	assertNotNil("function", function)
	q := r.Query.Select("withConstructor")
	q = q.Arg("function", function)

	return &TypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// TypeDefWithFieldOpts contains options for TypeDef.WithField
type TypeDefWithFieldOpts struct {
	// A doc string for the field, if any
	Description string
}

// Adds a static field for an Object TypeDef, failing if the type is not an object.
func (r *TypeDef) WithField(name string, typeDef *TypeDef, opts ...TypeDefWithFieldOpts) *TypeDef {
	assertNotNil("typeDef", typeDef)
	q := r.Query.Select("withField")
	for i := len(opts) - 1; i >= 0; i-- {
		// `description` optional argument
		if !querybuilder.IsZeroValue(opts[i].Description) {
			q = q.Arg("description", opts[i].Description)
		}
	}
	q = q.Arg("name", name)
	q = q.Arg("typeDef", typeDef)

	return &TypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// Adds a function for an Object or Interface TypeDef, failing if the type is not one of those kinds.
func (r *TypeDef) WithFunction(function *Function) *TypeDef {
	assertNotNil("function", function)
	q := r.Query.Select("withFunction")
	q = q.Arg("function", function)

	return &TypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// TypeDefWithInterfaceOpts contains options for TypeDef.WithInterface
type TypeDefWithInterfaceOpts struct {
	Description string
}

// Returns a TypeDef of kind Interface with the provided name.
func (r *TypeDef) WithInterface(name string, opts ...TypeDefWithInterfaceOpts) *TypeDef {
	q := r.Query.Select("withInterface")
	for i := len(opts) - 1; i >= 0; i-- {
		// `description` optional argument
		if !querybuilder.IsZeroValue(opts[i].Description) {
			q = q.Arg("description", opts[i].Description)
		}
	}
	q = q.Arg("name", name)

	return &TypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// Sets the kind of the type.
func (r *TypeDef) WithKind(kind TypeDefKind) *TypeDef {
	q := r.Query.Select("withKind")
	q = q.Arg("kind", kind)

	return &TypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// Returns a TypeDef of kind List with the provided type for its elements.
func (r *TypeDef) WithListOf(elementType *TypeDef) *TypeDef {
	assertNotNil("elementType", elementType)
	q := r.Query.Select("withListOf")
	q = q.Arg("elementType", elementType)

	return &TypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// TypeDefWithObjectOpts contains options for TypeDef.WithObject
type TypeDefWithObjectOpts struct {
	Description string
}

// Returns a TypeDef of kind Object with the provided name.
//
// Note that an object's fields and functions may be omitted if the intent is only to refer to an object. This is how functions are able to return their own object, or any other circular reference.
func (r *TypeDef) WithObject(name string, opts ...TypeDefWithObjectOpts) *TypeDef {
	q := r.Query.Select("withObject")
	for i := len(opts) - 1; i >= 0; i-- {
		// `description` optional argument
		if !querybuilder.IsZeroValue(opts[i].Description) {
			q = q.Arg("description", opts[i].Description)
		}
	}
	q = q.Arg("name", name)

	return &TypeDef{
		Query:  q,
		Client: r.Client,
	}
}

// Sets whether this type can be set to null.
func (r *TypeDef) WithOptional(optional bool) *TypeDef {
	q := r.Query.Select("withOptional")
	q = q.Arg("optional", optional)

	return &TypeDef{
		Query:  q,
		Client: r.Client,
	}
}

type CacheSharingMode string

func (CacheSharingMode) IsEnum() {}

const (
	// Shares the cache volume amongst many build pipelines, but will serialize the writes
	Locked CacheSharingMode = "LOCKED"

	// Keeps a cache volume for a single build pipeline
	Private CacheSharingMode = "PRIVATE"

	// Shares the cache volume amongst many build pipelines
	Shared CacheSharingMode = "SHARED"
)

type ImageLayerCompression string

func (ImageLayerCompression) IsEnum() {}

const (
	Estargz ImageLayerCompression = "EStarGZ"

	Gzip ImageLayerCompression = "Gzip"

	Uncompressed ImageLayerCompression = "Uncompressed"

	Zstd ImageLayerCompression = "Zstd"
)

type ImageMediaTypes string

func (ImageMediaTypes) IsEnum() {}

const (
	Dockermediatypes ImageMediaTypes = "DockerMediaTypes"

	Ocimediatypes ImageMediaTypes = "OCIMediaTypes"
)

type ModuleSourceKind string

func (ModuleSourceKind) IsEnum() {}

const (
	GitSource ModuleSourceKind = "GIT_SOURCE"

	LocalSource ModuleSourceKind = "LOCAL_SOURCE"
)

type NetworkProtocol string

func (NetworkProtocol) IsEnum() {}

const (
	Tcp NetworkProtocol = "TCP"

	Udp NetworkProtocol = "UDP"
)

type TypeDefKind string

func (TypeDefKind) IsEnum() {}

const (
	// A boolean value.
	BooleanKind TypeDefKind = "BOOLEAN_KIND"

	// A graphql input type, used only when representing the core API via TypeDefs.
	InputKind TypeDefKind = "INPUT_KIND"

	// An integer value.
	IntegerKind TypeDefKind = "INTEGER_KIND"

	// A named type of functions that can be matched+implemented by other objects+interfaces.
	//
	// Always paired with an InterfaceTypeDef.
	InterfaceKind TypeDefKind = "INTERFACE_KIND"

	// A list of values all having the same type.
	//
	// Always paired with a ListTypeDef.
	ListKind TypeDefKind = "LIST_KIND"

	// A named type defined in the GraphQL schema, with fields and functions.
	//
	// Always paired with an ObjectTypeDef.
	ObjectKind TypeDefKind = "OBJECT_KIND"

	// A string value.
	StringKind TypeDefKind = "STRING_KIND"

	// A special kind used to signify that no value is returned.
	//
	// This is used for functions that have no return value. The outer TypeDef specifying this Kind is always Optional, as the Void is never actually represented.
	VoidKind TypeDefKind = "VOID_KIND"
)
