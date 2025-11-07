package main

import (
	"fmt"
	"log"

	"github.com/7836246/kanggo"
)

func main() {
	fmt.Println("🚀 KangGo 请求验证系统演示")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	app := kanggo.Default()

	// 演示 1: 基础验证
	demo1BasicValidation(app)

	// 演示 2: 用户注册验证
	demo2UserRegistration(app)

	// 演示 3: 商品创建验证
	demo3ProductCreation(app)

	// 演示 4: 自定义验证规则
	demo4CustomValidation(app)

	// 首页 - HTML 表单
	app.GET("/", func(ctx *kanggo.Context) error {
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>KangGo 验证器演示</title>
    <style>
        body { font-family: Arial; max-width: 800px; margin: 50px auto; }
        .demo { border: 1px solid #ddd; padding: 20px; margin: 20px 0; }
        input, textarea { width: 100%; padding: 10px; margin: 5px 0; }
        button { padding: 10px 20px; margin: 10px 0; cursor: pointer; background: #4CAF50; color: white; border: none; }
        .error { color: red; margin: 10px 0; }
        .success { color: green; margin: 10px 0; }
        pre { background: #f5f5f5; padding: 10px; overflow-x: auto; }
    </style>
</head>
<body>
    <h1>🚀 KangGo 请求验证演示</h1>
    
    <div class="demo">
        <h2>演示 1: 基础验证</h2>
        <input type="text" id="basic-name" placeholder="姓名（必填）">
        <input type="email" id="basic-email" placeholder="邮箱">
        <input type="number" id="basic-age" placeholder="年龄（18-100）">
        <button onclick="testBasicValidation()">提交</button>
        <div id="basic-result"></div>
    </div>
    
    <div class="demo">
        <h2>演示 2: 用户注册</h2>
        <input type="text" id="reg-username" placeholder="用户名（3-20字符）">
        <input type="email" id="reg-email" placeholder="邮箱">
        <input type="password" id="reg-password" placeholder="密码（最少8位）">
        <button onclick="testUserRegistration()">注册</button>
        <div id="reg-result"></div>
    </div>
    
    <div class="demo">
        <h2>演示 3: 商品创建</h2>
        <input type="text" id="prod-name" placeholder="商品名称">
        <input type="number" id="prod-price" placeholder="价格（0.01-999999）">
        <input type="number" id="prod-stock" placeholder="库存">
        <textarea id="prod-desc" placeholder="描述（最少10字符）" rows="3"></textarea>
        <button onclick="testProductCreation()">创建</button>
        <div id="prod-result"></div>
    </div>

    <script>
        function testBasicValidation() {
            const data = {
                name: document.getElementById('basic-name').value,
                email: document.getElementById('basic-email').value,
                age: parseInt(document.getElementById('basic-age').value) || 0
            };
            
            fetch('/api/validate/basic', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            })
            .then(r => r.json())
            .then(showResult('basic-result'))
            .catch(showError('basic-result'));
        }

        function testUserRegistration() {
            const data = {
                username: document.getElementById('reg-username').value,
                email: document.getElementById('reg-email').value,
                password: document.getElementById('reg-password').value
            };
            
            fetch('/api/users/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            })
            .then(r => r.json())
            .then(showResult('reg-result'))
            .catch(showError('reg-result'));
        }

        function testProductCreation() {
            const data = {
                name: document.getElementById('prod-name').value,
                price: parseFloat(document.getElementById('prod-price').value) || 0,
                stock: parseInt(document.getElementById('prod-stock').value) || 0,
                description: document.getElementById('prod-desc').value
            };
            
            fetch('/api/products', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            })
            .then(r => r.json())
            .then(showResult('prod-result'))
            .catch(showError('prod-result'));
        }

        function showResult(id) {
            return function(data) {
                const div = document.getElementById(id);
                if (data.error) {
                    div.innerHTML = '<div class="error">❌ ' + data.error + '</div>';
                } else {
                    div.innerHTML = '<div class="success">✅ 验证成功！</div><pre>' + JSON.stringify(data, null, 2) + '</pre>';
                }
            };
        }

        function showError(id) {
            return function(err) {
                document.getElementById(id).innerHTML = '<div class="error">❌ ' + err.message + '</div>';
            };
        }
    </script>
</body>
</html>
`
		return ctx.HTML(200, html)
	})

	fmt.Println("🌐 服务器启动在 http://localhost:8080")
	fmt.Println()
	fmt.Println("测试地址：")
	fmt.Println("  • 基础验证:  POST /api/validate/basic")
	fmt.Println("  • 用户注册:  POST /api/users/register")
	fmt.Println("  • 商品创建:  POST /api/products")
	fmt.Println("  • 网页演示:  http://localhost:8080")
	fmt.Println()

	app.Run(":8080")
}

// demo1BasicValidation 演示 1: 基础验证
func demo1BasicValidation(app *kanggo.KangGo) {
	type BasicRequest struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"email"`
		Age   int    `json:"age" validate:"min=18,max=100"`
	}

	app.POST("/api/validate/basic", func(ctx *kanggo.Context) error {
		var req BasicRequest

		// 使用 ValidateJSON 自动验证
		if err := ctx.ValidateJSON(&req); err != nil {
			log.Printf("验证失败: %v", err)
			return ctx.JSON(400, map[string]string{
				"error": err.Error(),
			})
		}

		log.Printf("验证成功: %+v", req)

		return ctx.JSON(200, map[string]interface{}{
			"message": "验证成功",
			"data":    req,
		})
	})
}

// demo2UserRegistration 演示 2: 用户注册验证
func demo2UserRegistration(app *kanggo.KangGo) {
	type RegisterRequest struct {
		Username string `json:"username" validate:"required,minlen=3,maxlen=20,alphanum"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,minlen=8"`
	}

	app.POST("/api/users/register", func(ctx *kanggo.Context) error {
		var req RegisterRequest

		if err := ctx.ValidateJSON(&req); err != nil {
			// 返回详细的验证错误
			if validationErrs, ok := err.(kanggo.ValidationErrors); ok {
				errors := make([]map[string]string, 0)
				for _, e := range validationErrs {
					errors = append(errors, map[string]string{
						"field":   e.Field,
						"message": e.Message,
						"tag":     e.Tag,
					})
				}
				return ctx.JSON(400, map[string]interface{}{
					"error":  "验证失败",
					"errors": errors,
				})
			}

			return ctx.JSON(400, map[string]string{
				"error": err.Error(),
			})
		}

		// 验证成功，创建用户
		log.Printf("用户注册: %s (%s)", req.Username, req.Email)

		return ctx.JSON(201, map[string]interface{}{
			"message": "注册成功",
			"user": map[string]interface{}{
				"username": req.Username,
				"email":    req.Email,
				"id":       12345,
			},
		})
	})
}

// demo3ProductCreation 演示 3: 商品创建验证
func demo3ProductCreation(app *kanggo.KangGo) {
	type ProductRequest struct {
		Name        string  `json:"name" validate:"required,minlen=2,maxlen=100"`
		Price       float64 `json:"price" validate:"required,min=0.01,max=999999"`
		Stock       int     `json:"stock" validate:"required,min=0"`
		Description string  `json:"description" validate:"minlen=10,maxlen=1000"`
		SKU         string  `json:"sku" validate:"alphanum"`
	}

	app.POST("/api/products", func(ctx *kanggo.Context) error {
		var req ProductRequest

		if err := ctx.ValidateJSON(&req); err != nil {
			log.Printf("商品验证失败: %v", err)
			return ctx.JSON(400, map[string]interface{}{
				"error":   "验证失败",
				"message": err.Error(),
			})
		}

		log.Printf("商品创建: %s - $%.2f", req.Name, req.Price)

		return ctx.JSON(201, map[string]interface{}{
			"message": "商品创建成功",
			"product": map[string]interface{}{
				"id":          12345,
				"name":        req.Name,
				"price":       req.Price,
				"stock":       req.Stock,
				"description": req.Description,
			},
		})
	})
}

// demo4CustomValidation 演示 4: 自定义验证
func demo4CustomValidation(app *kanggo.KangGo) {
	type CustomRequest struct {
		Username string `json:"username" validate:"required,regex=^[a-z][a-z0-9_]{2,19}$"`
		Website  string `json:"website" validate:"url"`
		Code     string `json:"code" validate:"len=6,numeric"`
	}

	app.POST("/api/validate/custom", func(ctx *kanggo.Context) error {
		var req CustomRequest

		if err := ctx.ValidateJSON(&req); err != nil {
			return ctx.JSON(400, map[string]interface{}{
				"error":   "自定义验证失败",
				"message": err.Error(),
			})
		}

		return ctx.JSON(200, map[string]interface{}{
			"message": "自定义验证成功",
			"data":    req,
		})
	})
}
