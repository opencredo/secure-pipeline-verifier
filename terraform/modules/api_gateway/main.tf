resource "aws_api_gateway_resource" "resource" {
  path_part   = var.path_part
  parent_id   = var.root_resource_id
  rest_api_id = var.api_id
}

resource "aws_api_gateway_method" "method" {
  rest_api_id   = var.api_id
  resource_id   = aws_api_gateway_resource.resource.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_method_response" "response_200" {
  rest_api_id = var.api_id
  resource_id = aws_api_gateway_resource.resource.id
  http_method = aws_api_gateway_method.method.http_method
  status_code = "200"
}

resource "aws_api_gateway_integration_response" "integration_response" {
  rest_api_id = var.api_id
  resource_id = aws_api_gateway_resource.resource.id
  http_method = aws_api_gateway_method.method.http_method
  status_code = aws_api_gateway_method_response.response_200.status_code

  depends_on = [
    aws_api_gateway_integration.lambda_integration
  ]

}

resource "aws_api_gateway_integration" "lambda_integration" {
  rest_api_id             = var.api_id
  resource_id             = aws_api_gateway_resource.resource.id
  http_method             = aws_api_gateway_method.method.http_method
  integration_http_method = aws_api_gateway_method.method.http_method
  type                    = "AWS"
  #  . For AWS integrations, the URI should be of the form
  #  arn:aws:apigateway:{region}:{subdomain.service|service}:{path|action}/{service_api}
  uri                     = var.invoke_arn

  passthrough_behavior = var.passthrough_behavior
  request_templates = {
    # If the POST request is not JSON, then map it for Lambda.
    "application/x-www-form-urlencoded" = var.urlencoded_tmpl
  }
}

output "lambda_integration" {
  description = "Deployment invoke url"
  value       = aws_api_gateway_integration.lambda_integration
}

# Permissions
resource "aws_lambda_permission" "allow_api_gw_to_call_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = var.function_name
  principal     = "apigateway.amazonaws.com"
  # More: http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-control-access-using-iam-policies-to-invoke-api.html
  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${var.api_id}/*/${aws_api_gateway_method.method.http_method}${aws_api_gateway_resource.resource.path}"
}