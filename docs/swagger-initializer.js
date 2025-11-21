window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">

 // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    // 👇 CRITICAL CHANGE: This tells Swagger UI to look for your specification file 
    // in the same directory (the /docs folder).
    url: "./openapi.yaml", 
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });

//</editor-fold>
};