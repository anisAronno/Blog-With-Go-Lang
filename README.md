# Go Web App - Laravel-Inspired Blog Platform

A comprehensive Go web application with Laravel-like structure, featuring CRUD operations, user authentication, and modern Tailwind CSS design. Perfect for learning Go web development with familiar patterns.

## 🚀 Features

- **User Authentication** - Login, register, logout with sessions
- **Blog CRUD Operations** - Create, read, update, delete blog posts
- **Modern Dashboard** - Beautiful Tailwind CSS interface
- **Database Integration** - MySQL with proper migrations and seeders
- **Security Features** - Password hashing, CSRF protection, input validation
- **Responsive Design** - Mobile-friendly interface
- **Unit Testing** - Comprehensive test suite
- **Enterprise Structure** - Scalable, maintainable codebase

## 📁 Project Structure

```
├── app/
│   ├── controllers/        # HTTP request handlers (like Laravel controllers)
│   │   ├── auth_controller.go      # Authentication logic
│   │   ├── blog_controller.go      # Blog CRUD operations
│   │   ├── dashboard_controller.go # Dashboard pages
│   │   ├── home_controller.go      # Public pages
│   │   └── controller.go           # Base controller utilities
│   ├── models/            # Database models (like Laravel models)
│   │   ├── user.go        # User model with authentication
│   │   └── blog.go        # Blog model with CRUD operations
│   └── middleware/        # HTTP middleware (like Laravel middleware)
│       └── auth.go        # Authentication and session middleware
├── config/                # Configuration management
│   └── config.go         # Database and app configuration
├── database/
│   ├── migrations/       # Database schema migrations
│   │   └── migrate.go
│   └── seeders/         # Database seeders for test data
│       └── seed.go
├── public/              # Static assets
│   ├── css/
│   │   └── styles.css   # Custom CSS styles
│   └── js/
│       └── app.js       # JavaScript functionality
├── routes/              # Route definitions (like Laravel routes)
│   └── web.go          # Web routes configuration
├── templates/           # HTML templates (like Laravel views)
│   ├── layouts/
│   │   └── base.html    # Base layout template
│   ├── auth/            # Authentication pages
│   ├── dashboard/       # Dashboard pages
│   └── blog/           # Blog-related pages
├── tests/              # Unit tests
│   ├── models_test.go  # Model tests
│   └── controllers_test.go # Controller tests
├── .env               # Environment configuration
├── go.mod            # Go module definition
└── main.go          # Application entry point
```

## 🛠️ Getting Started

### Prerequisites

- **Go 1.21+** - [Download & Install Go](https://golang.org/dl/)
- **MySQL 8.0+** - [Download & Install MySQL](https://dev.mysql.com/downloads/)
- **Git** - For cloning the repository

### Installation

1. **Clone the Repository**

   ```bash
   git clone <repository-url>
   cd go-web-app
   ```

2. **Install Dependencies**

   ```bash
   go mod tidy
   ```

3. **Setup Database**

   ```bash
   # Create MySQL database
   mysql -u root -p
   CREATE DATABASE go_web_app;
   CREATE DATABASE go_web_app_test; -- For testing
   EXIT;
   ```

4. **Configure Environment**

   ```bash
   # Copy environment file
   cp .env.example .env

   # Edit .env file with your database credentials
   nano .env
   ```

   Update the `.env` file:

   ```env
   # Database Configuration
   DB_HOST=localhost
   DB_PORT=3306
   DB_NAME=go_web_app
   DB_USER=root
   DB_PASSWORD=your_mysql_password

   # Application Configuration
   APP_PORT=3000
   APP_ENV=development
   APP_KEY=your-secret-key-here

   # Session Configuration
   SESSION_SECRET=your-session-secret-here
   ```

5. **Run Database Migrations**

   ```bash
   go run database/migrations/migrate.go
   ```

6. **Seed Database with Test Data**

   ```bash
   go run database/seeders/seed.go
   ```

7. **Start the Application**

   ```bash
   go run main.go
   ```

8. **Access the Application**
   - **Homepage:** http://localhost:3000
   - **Login:** http://localhost:3000/login
   - **Register:** http://localhost:3000/register
   - **Dashboard:** http://localhost:3000/dashboard (after login)

## 👤 Default Test Users

After running the seeder, you can login with these accounts:

| Email             | Password | Role         |
| ----------------- | -------- | ------------ |
| admin@example.com | password | Admin User   |
| john@example.com  | password | Regular User |
| jane@example.com  | password | Regular User |

## 🧪 Running Tests

```bash
# Run all tests
go test ./tests/...

# Run specific test file
go test ./tests/models_test.go

# Run tests with verbose output
go test -v ./tests/...

# Run tests with coverage
go test -cover ./tests/...
```

## 📚 API Endpoints

### Public Routes

- `GET /` - Homepage with blog listing
- `GET /blog/{id}` - View individual blog post
- `GET /login` - Login page
- `POST /login` - Process login
- `GET /register` - Registration page
- `POST /register` - Process registration

### Protected Routes (Require Authentication)

- `GET /dashboard` - Main dashboard
- `GET /dashboard/profile` - User profile
- `GET /dashboard/users` - All users listing
- `GET /dashboard/blogs` - User's blog management
- `GET /dashboard/blogs/create` - Create new blog form
- `POST /dashboard/blogs` - Store new blog
- `GET /dashboard/blogs/{id}/edit` - Edit blog form
- `POST /dashboard/blogs/{id}` - Update blog
- `POST /dashboard/blogs/{id}/delete` - Delete blog
- `POST /logout` - Logout user

## 🏗️ Architecture & Design Patterns

### MVC Architecture

- **Models** (`app/models/`) - Data layer with database operations
- **Views** (`templates/`) - Presentation layer with HTML templates
- **Controllers** (`app/controllers/`) - Business logic and request handling

### Key Design Patterns

- **Repository Pattern** - Models act as repositories for data access
- **Middleware Pattern** - Authentication, logging, CORS handling
- **Template Pattern** - Consistent HTML layout inheritance
- **Dependency Injection** - Loose coupling between components

### Security Features

- **Password Hashing** - bcrypt for secure password storage
- **Session Management** - Secure session handling with gorilla/sessions
- **Input Validation** - Server-side validation for all forms
- **SQL Injection Prevention** - Prepared statements for database queries
- **XSS Protection** - Template auto-escaping

## 🎨 Frontend Technologies

- **Tailwind CSS** - Utility-first CSS framework
- **Font Awesome** - Icon library
- **Vanilla JavaScript** - Modern ES6+ features
- **Responsive Design** - Mobile-first approach

## 🔧 Development Tips

### Adding New Features

1. **Create Model** (if needed)

   ```go
   // app/models/your_model.go
   type YourModel struct {
       ID        int       `json:"id"`
       Name      string    `json:"name"`
       CreatedAt time.Time `json:"created_at"`
   }
   ```

2. **Create Controller**

   ```go
   // app/controllers/your_controller.go
   func (c *YourController) Index(w http.ResponseWriter, r *http.Request) {
       // Handle request
   }
   ```

3. **Add Routes**

   ```go
   // routes/web.go
   r.HandleFunc("/your-route", yourController.Index).Methods("GET")
   ```

4. **Create Templates**
   ```html
   <!-- templates/your_template.html -->
   {{define "content"}}
   <!-- Your HTML content -->
   {{end}}
   ```

### Database Migrations

Create new migration files in `database/migrations/` following the pattern:

```go
// Add new table
CREATE TABLE IF NOT EXISTS your_table (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Environment Configuration

Add new configuration options in `config/config.go`:

```go
type Config struct {
    // Existing fields...
    YourNewSetting string
}

// In LoadConfig()
YourNewSetting: getEnv("YOUR_NEW_SETTING", "default_value"),
```

## 🚀 Deployment

### Production Checklist

1. **Environment Setup**

   - Set `APP_ENV=production`
   - Use strong, unique `APP_KEY` and `SESSION_SECRET`
   - Configure production database

2. **Security**

   - Enable HTTPS
   - Set secure session options
   - Configure firewall rules

3. **Performance**
   - Enable Go build optimizations
   - Configure database connection pooling
   - Set up static file serving (nginx/Apache)

### Docker Deployment (Optional)

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/public ./public
CMD ["./main"]
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 Learning Resources

### Go Web Development

- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Web Programming](https://github.com/astaxie/build-web-application-with-golang)

### Similar Patterns in Other Languages

- **Laravel (PHP)** - Similar MVC structure and conventions
- **Django (Python)** - Models, views, templates pattern
- **Ruby on Rails** - Convention over configuration

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Inspired by Laravel's elegant syntax and structure
- Built with Go's powerful standard library
- Styled with Tailwind CSS for modern design
- Icons provided by Font Awesome

---

**Happy Coding! 🎉**

This project demonstrates enterprise-grade Go web development with familiar Laravel patterns, making it perfect for learning and real-world applications.
