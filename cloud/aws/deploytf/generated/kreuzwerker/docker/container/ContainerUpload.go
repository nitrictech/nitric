package container


type ContainerUpload struct {
	// Path to the file in the container where is upload goes to.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#file Container#file}
	File *string `field:"required" json:"file" yaml:"file"`
	// Literal string value to use as the object content, which will be uploaded as UTF-8-encoded text.
	//
	// Conflicts with `content_base64` & `source`
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#content Container#content}
	Content *string `field:"optional" json:"content" yaml:"content"`
	// Base64-encoded data that will be decoded and uploaded as raw bytes for the object content.
	//
	// This allows safely uploading non-UTF8 binary data, but is recommended only for larger binary content such as the result of the `base64encode` interpolation function. See [here](https://github.com/terraform-providers/terraform-provider-docker/issues/48#issuecomment-374174588) for the reason. Conflicts with `content` & `source`
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#content_base64 Container#content_base64}
	ContentBase64 *string `field:"optional" json:"contentBase64" yaml:"contentBase64"`
	// If `true`, the file will be uploaded with user executable permission. Defaults to `false`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#executable Container#executable}
	Executable interface{} `field:"optional" json:"executable" yaml:"executable"`
	// A filename that references a file which will be uploaded as the object content.
	//
	// This allows for large file uploads that do not get stored in state. Conflicts with `content` & `content_base64`
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#source Container#source}
	Source *string `field:"optional" json:"source" yaml:"source"`
	// If using `source`, this will force an update if the file content has updated but the filename has not.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#source_hash Container#source_hash}
	SourceHash *string `field:"optional" json:"sourceHash" yaml:"sourceHash"`
}

