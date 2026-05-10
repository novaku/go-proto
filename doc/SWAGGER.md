# Swagger Documentation Quick Reference

## Overview
This project automatically generates OpenAPI/Swagger documentation from your Protocol Buffer definitions. The documentation is generated every time you run `make generate`.

## Files Generated
- **Swagger JSON**: `gen/guestbook/v1/guestbook.swagger.json`
  - OpenAPI 2.0 specification
  - Contains all API endpoints, schemas, and documentation
  - Can be imported into any OpenAPI-compatible tool

## Viewing the Documentation

### Method 1: Swagger UI (Recommended)
1. Run the Swagger server:
   ```bash
   make swagger
   ```
2. Open your browser to: `http://localhost:8081/swagger-ui.html`
3. You'll see an interactive API documentation with:
   - All endpoints listed
   - Request/response schemas
   - Try-it-out functionality
   - Field descriptions

### Method 2: Raw JSON
View the raw Swagger specification:
```bash
cat gen/guestbook/v1/guestbook.swagger.json
```

### Method 3: Import to External Tools
You can import the `guestbook.swagger.json` file into:
- **Postman**: File → Import → Upload the JSON file
- **Swagger Editor**: https://editor.swagger.io/ → File → Import file
- **Insomnia**: Import → From file

## Adding Documentation to Your Proto Files

### Service-level Documentation
Add descriptions to your RPC methods:
```protobuf
rpc AddEntry(AddEntryRequest) returns (AddEntryResponse) {
  option (google.api.http) = {
    post: "/v1/guestbook"
    body: "*"
  };
  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
    summary: "Add guestbook entry";
    description: "Creates a new entry in the guestbook with name, email, and message";
    tags: "Guestbook";
  };
}
```

### Field-level Documentation
Add comments to message fields:
```protobuf
message AddEntryRequest {
  string name = 1;     // Name of the person signing the guestbook
  string email = 2;    // Email address of the person
  string message = 3;  // Message to be added to the guestbook
}
```

### API-level Metadata
Configure global API information:
```protobuf
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  info: {
    title: "Guestbook API";
    version: "1.0";
    description: "A simple guestbook API built with gRPC and grpc-gateway";
    contact: {
      name: "Guestbook API Support";
      email: "support@example.com";
    };
  };
  schemes: HTTP;
  schemes: HTTPS;
  consumes: "application/json";
  produces: "application/json";
};
```

## Regenerating Documentation

Whenever you modify your `.proto` files, regenerate the documentation:
```bash
make generate
```

This will:
1. Install/update required tools
2. Download proto dependencies
3. Generate Go code
4. Generate updated Swagger documentation

## Customization Options

The Swagger generation supports many options. Common ones include:

- `logtostderr=true` - Log errors to stderr
- `allow_merge=true` - Merge multiple proto files
- `json_names_for_fields=true` - Use JSON naming for fields
- `simple_operation_ids=true` - Simplify operation IDs

To add options, modify the `--openapiv2_opt` flag in the Makefile:
```makefile
--openapiv2_out=gen --openapiv2_opt=logtostderr=true,simple_operation_ids=true \
```

## Troubleshooting

### Proto dependencies not found
If you see errors about missing proto files, run:
```bash
make proto-deps
```

### Swagger UI not loading
1. Ensure you've run `make generate` first
2. Check that `gen/guestbook/v1/guestbook.swagger.json` exists
3. Verify Python 3 is installed: `python3 --version`

### Port 8081 already in use
Change the port in the Makefile's `swagger` target:
```makefile
swagger:
	@cd gen && python3 -m http.server 8082
```

## Additional Resources

- [grpc-gateway Documentation](https://grpc-ecosystem.github.io/grpc-gateway/)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Protocol Buffers Guide](https://protobuf.dev/)
