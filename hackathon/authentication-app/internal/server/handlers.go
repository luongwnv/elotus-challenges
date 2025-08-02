package server

import (
	"authentication-app/internal/controllers"
	"authentication-app/internal/middleware"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	golog "github.com/luongwnv/go-log"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func HandlePanic(contextName string) {
	if r := recover(); r != nil {
		fmt.Println("An error occurred:", r)
	}
}

// @title SIMPLE AUTHENTICATION APP API
// @version 1.0
// @description API documentation for SIMPLE AUTHENTICATION APP services
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath /
// @security BearerAuth
func (s *Server) MapHandlers() error {
	app := s.fiber
	// validate := validator.New()

	app.Use(func(c *fiber.Ctx) error {
		return c.Next()
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(func(c *fiber.Ctx) error {
		defer HandlePanic("HTTP Service")
		return c.Next()
	})

	app.Get("/api/swagger/*", fiberSwagger.WrapHandler)

	app.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02T15:04:05.999",
		Format:     "[${ip}] ${time} ${locals:requestid} ${method} ${path} ${status} ${latency}\n",
	}))

	// HTML pages
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Type("html").SendString(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Authentication App</title>
				<style>
					body { font-family: Arial, sans-serif; max-width: 600px; margin: 50px auto; padding: 20px; }
					.container { background: #f5f5f5; padding: 30px; border-radius: 8px; }
					input, button { padding: 10px; margin: 5px 0; width: 100%; box-sizing: border-box; }
					button { background: #007bff; color: white; border: none; cursor: pointer; }
					button:hover { background: #0056b3; }
					.hidden { display: none; }
					.success { color: green; }
					.error { color: red; }
				</style>
			</head>
			<body>
				<div class="container">
					<h2>Login</h2>
					<div id="loginForm">
						<input type="text" id="username" placeholder="Username" required>
						<input type="password" id="password" placeholder="Password" required>
						<button onclick="login()">Login</button>
						<button onclick="register()">Register</button>
					</div>
					
					<div id="uploadForm" class="hidden">
						<h3>Upload File</h3>
						<input type="file" id="fileInput">
						<button onclick="uploadFile()">Upload</button>
						<button onclick="logout()">Logout</button>
					</div>
					
					<div id="message"></div>
				</div>

				<script>
					let token = localStorage.getItem('authToken');
					
					if (token) {
						showUploadForm();
					}

					function showMessage(msg, isError = false) {
						const messageDiv = document.getElementById('message');
						messageDiv.textContent = msg;
						messageDiv.className = isError ? 'error' : 'success';
					}

					function showUploadForm() {
						document.getElementById('loginForm').classList.add('hidden');
						document.getElementById('uploadForm').classList.remove('hidden');
					}

					function showLoginForm() {
						document.getElementById('loginForm').classList.remove('hidden');
						document.getElementById('uploadForm').classList.add('hidden');
					}

					async function register() {
						const username = document.getElementById('username').value;
						const password = document.getElementById('password').value;

						try {
							const response = await fetch('/auth/register', {
								method: 'POST',
								headers: { 'Content-Type': 'application/json' },
								body: JSON.stringify({ username, password })
							});

							const data = await response.json();
							if (response.ok) {
								localStorage.setItem('authToken', data.token);
								token = data.token;
								showUploadForm();
								showMessage('Registration successful!');
							} else {
								showMessage(data.error, true);
							}
						} catch (error) {
							showMessage('Network error', true);
						}
					}

					async function login() {
						const username = document.getElementById('username').value;
						const password = document.getElementById('password').value;

						try {
							const response = await fetch('/auth/login', {
								method: 'POST',
								headers: { 'Content-Type': 'application/json' },
								body: JSON.stringify({ username, password })
							});

							const data = await response.json();
							if (response.ok) {
								localStorage.setItem('authToken', data.token);
								token = data.token;
								showUploadForm();
								showMessage('Login successful!');
							} else {
								showMessage(data.error, true);
							}
						} catch (error) {
							showMessage('Network error', true);
						}
					}

					async function uploadFile() {
						const fileInput = document.getElementById('fileInput');
						const file = fileInput.files[0];

						if (!file) {
							showMessage('Please select a file', true);
							return;
						}

						const formData = new FormData();
						formData.append('file', file);

						try {
							const response = await fetch('/files/upload', {
								method: 'POST',
								headers: { 'Authorization': 'Bearer ' + token },
								body: formData
							});

							const data = await response.json();
							if (response.ok) {
								showMessage('File uploaded successfully!');
								fileInput.value = '';
							} else {
								showMessage(data.error, true);
							}
						} catch (error) {
							showMessage('Upload failed', true);
						}
					}

					function logout() {
						localStorage.removeItem('authToken');
						token = null;
						showLoginForm();
						showMessage('Logged out successfully!');
					}
				</script>
			</body>
			</html>
		`)
	})

	// Health check routes
	monitoringHandler := controllers.NewHandler(s.cfg, s.logger, s.rdbIns)
	app.Get("/api/readiness", timeout.New(monitoringHandler.Readiness, time.Duration(s.cfg.ServerCtxDefaultTimeout)*time.Second))
	app.Get("/api/liveness", monitoringHandler.Liveness)

	// Auth routes
	authController := controllers.NewAuthController(s.cfg, s.logger, s.rdbIns)
	authGroup := app.Group("/auth")
	authGroup.Post("/register", authController.Register)
	authGroup.Post("/login", authController.Login)

	jwtMiddleware := middleware.JWTAuth(s.cfg, s.rdbIns)
	authGroup.Post("/revoke", jwtMiddleware, authController.RevokeToken)

	fileController := controllers.NewFileController(s.logger, s.rdbIns)
	fileGroup := app.Group("/files")
	fileGroup.Post("/upload", jwtMiddleware, fileController.UploadFile)

	golog.Info("Loaded all route!")

	return nil
}
