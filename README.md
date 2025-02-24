
# API Documentation & Workflow Guide

# Secure Your Web Assets and APIs with Kraken Integration

  
  

 1. **Authentication**
 - **Signup**
 Endpoint: POST /signup
 Request:
```json
{
"username": "your_username",
"password": "your_password"
}
```
- **Login**
Endpoint: POST /login
Request: Same as signup.
Response:
```json
{"token": "your_jwt_token"}
```
Note: Include this token in the Authorization: Bearer <token> header for all protected routes.

  

2. **URL Security Flow**
Secure web URLs by discovering vulnerabilities with ZAP.
Step 1: Add URLs
Endpoint: POST /api/v1/add/url
Request:
```json
{
"service": "my-webapp",
"url_list": "https://example.com,https://test.com"
}
```
- Start Spider Scan (Discovery)
Endpoint: POST /api/v1/url-discovery/service/my-webapp
Response:
```json
{"message": "Spider scan initiated"}
```
- Start Active Scan
Endpoint: POST /api/v1/url-scan/service/my-webapp
Response:
```json
{"message": "Active scan initiated"}
```
- Check Security Status
Endpoint: GET /api/v1/url-report/service/my-webapp
Response:
```json
{"message": "Your project is safe. Proceed."}
// OR
{"message": "Your project failed the security check."}
```
3. **API Security Flow**
Secure API schemas (e.g., OpenAPI) with ZAP.
- Add API Target
Endpoint: POST /api/v1/add/api
Request:
```json
{
"service": "my-api",
"api_schema": "https://api.example.com/swagger.json"
}
```
- Import API to ZAP
Endpoint: POST /api/v1/import-api/service/my-api
Response:
```json
{"message": "API target imported"}
```
- Start API Scan
Endpoint: POST /api/v1/api-scan/service/my-api
Response:
```json
{"message": "API scan started"}
```
- Retrieve API Alerts
Endpoint: GET /api/v1/api-report/service/my-api
Response:
```json
{
"message": "Your API is safe. Proceed.",
"counts": {"Medium": 0, "High": 0, "Critical": 0}
}
```
4. Example Commands
Run a Spider Scan:
```bash
curl -X POST http://localhost:8080/api/v1/url-discovery/service/my-webapp \
-H "Authorization: Bearer YOUR_TOKEN"
```
Check API Security:
```bash
curl http://localhost:8080/api/v1/api-report/service/my-api \
-H "Authorization: Bearer YOUR_TOKEN"
```
Key Notes
Authorization: All protected routes require a valid JWT token.
User Isolation: Data is scoped to the authenticated user.
ZAP Dependency: Scans will fail if ZAP is not running.
This is your roadmap to securing web and API assets. Follow the steps, run the commands, and ensure your systems are vulnerability-free! 🔒#kraken
