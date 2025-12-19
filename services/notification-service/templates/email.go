package templates

import (
	"bytes"
	"html/template"
)

// OrderConfirmationTemplate generates order confirmation HTML
var OrderConfirmationHTML = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Order Confirmation</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #4F46E5; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .order-item { border-bottom: 1px solid #ddd; padding: 10px 0; }
        .total { font-size: 18px; font-weight: bold; color: #4F46E5; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Order Confirmed!</h1>
        </div>
        <div class="content">
            <p>Hi {{.CustomerName}},</p>
            <p>Thank you for your order! We're excited to let you know that we've received your order #{{.OrderID}}.</p>
            
            <h3>Order Details:</h3>
            {{range .Items}}
            <div class="order-item">
                <strong>{{.ProductName}}</strong><br>
                Quantity: {{.Quantity}} × ${{printf "%.2f" .Price}}
            </div>
            {{end}}
            
            <p class="total">Total: ${{printf "%.2f" .Total}}</p>
            
            <p>We'll send you another email when your order ships.</p>
        </div>
        <div class="footer">
            <p>© 2024 Your E-Commerce Store. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

// WelcomeEmailHTML generates welcome email HTML
var WelcomeEmailHTML = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome!</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 40px; text-align: center; }
        .content { padding: 30px; background: #ffffff; }
        .button { display: inline-block; background: #4F46E5; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome, {{.Name}}! 🎉</h1>
        </div>
        <div class="content">
            <p>We're thrilled to have you join our community!</p>
            <p>Here's what you can do now:</p>
            <ul>
                <li>Browse our latest products</li>
                <li>Set up your preferences</li>
                <li>Enjoy exclusive member discounts</li>
            </ul>
            <p style="text-align: center;">
                <a href="{{.ShopURL}}" class="button">Start Shopping</a>
            </p>
        </div>
        <div class="footer">
            <p>© 2024 Your E-Commerce Store. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

type TemplateEngine struct {
	templates map[string]*template.Template
}

func NewTemplateEngine() *TemplateEngine {
	engine := &TemplateEngine{
		templates: make(map[string]*template.Template),
	}
	engine.loadTemplates()
	return engine
}

func (e *TemplateEngine) loadTemplates() {
	e.templates["order_confirmation"], _ = template.New("order_confirmation").Parse(OrderConfirmationHTML)
	e.templates["welcome"], _ = template.New("welcome").Parse(WelcomeEmailHTML)
}

func (e *TemplateEngine) Render(templateName string, data any) (string, error) {
	tmpl, ok := e.templates[templateName]
	if !ok {
		return "", nil
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
