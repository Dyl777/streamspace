# StreamSpace Website

This directory contains the static website for StreamSpace, hosted on GitHub Pages.

##  Structure

```
website/
├── index.html              # Homepage
├── features.html           # Features page
├── docs.html              # Documentation
├── getting-started.html   # Installation guide
├── plugins.html           # Plugin system showcase
├── templates.html         # Application templates showcase
├── css/
│   └── style.css          # Main stylesheet
├── js/
│   └── main.js            # JavaScript functionality
└── README.md              # This file
```

##  Deployment

### GitHub Pages Setup

1. **Enable GitHub Pages** in repository settings:
   - Go to Settings → Pages
   - Source: Deploy from a branch
   - Branch: `main` (or your branch)
   - Folder: `/website`
   - Save

2. **Access the site**:
   - URL: `https://joshuaaferguson.github.io/streamspace/`
   - Custom domain (optional): Configure CNAME file

### Local Development

Serve the website locally:

```bash
# Using Python
cd website
python3 -m http.server 8000

# Using Node.js (http-server)
npx http-server website -p 8000

# Using PHP
cd website
php -S localhost:8000
```

Then open `http://localhost:8000` in your browser.

##  Customization

### Colors

Edit `css/style.css` to change the color scheme:

```css
:root {
  --primary: #6366f1;        /* Primary color */
  --secondary: #8b5cf6;      /* Secondary color */
  --background: #0f172a;     /* Background */
  --surface: #1e293b;        /* Surface elements */
  --text: #f1f5f9;          /* Text color */
}
```

### Content

- **Homepage**: Edit `index.html`
- **Features**: Edit `features.html`
- **Documentation**: Edit `docs.html`
- **Getting Started**: Edit `getting-started.html`
- **Plugins**: Edit `plugins.html`
- **Templates**: Edit `templates.html`

### Navigation

Update navigation in all pages by editing the `<nav>` section in each HTML file.

## � Responsive Design

The website is fully responsive and works on:
- Desktop browsers
- Tablets
- Mobile phones

Mobile menu activates automatically on screens < 768px wide.

## ✨ Features

- **Modern Design**: Clean, professional design with dark theme
- **Responsive**: Works on all screen sizes
- **Fast**: Minimal JavaScript, optimized CSS
- **SEO Friendly**: Proper meta tags and semantic HTML
- **Accessible**: ARIA labels and keyboard navigation
- **Smooth Animations**: Fade-in effects and transitions

##  Technologies

- **HTML5**: Semantic markup
- **CSS3**: Modern features (Grid, Flexbox, Custom Properties)
- **JavaScript**: Vanilla JS (no frameworks)
- **Fonts**: Google Fonts (Inter)

##  License

MIT License - same as StreamSpace project.

## � Contributing

To improve the website:

1. Edit the HTML/CSS/JS files
2. Test locally
3. Submit a pull request

##  Support

- [GitHub Issues](https://github.com/JoshuaAFerguson/streamspace/issues)
- [Discussions](https://github.com/JoshuaAFerguson/streamspace/discussions)
