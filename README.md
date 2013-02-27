json-tree
=========

Parse and build JSON as tree data-structure with Go
---------------------------------------------------

json-tree allows you to dynamically query data from the map[string]interface{} based
tree data-structure generated by json.Unmarshal().

See:
	
	FromBytes(jsonBytes []byte) (Tree, error)
	FromString(jsonString string) (Tree, error)
	FromFile(filename string) (Tree, error)
	FromURL(url string) (Tree, error)
	FromReader(reader io.Reader) (Tree, error)


Builder creates such json.Marshal() compatible trees.

Example:

	var builder json.Builder
	builder.BeginObject().Name("Greeting").Value("Hello World!").EndObject()

	// Create tree and query "Greeting"
	greeting := builder.Tree().Select("Greeting").GetString()