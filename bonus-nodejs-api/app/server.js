const express = require('express');
const os = require('os');
const app = express();
const port = process.env.PORT || 3000;

app.get('/', (req, res) => {
  res.json({
    message: "🚀 Hello from the Node.js API running on AWS EKS!",
    hostname: os.hostname(),
    platform: os.platform(),
    timestamp: new Date().toISOString()
  });
});

app.get('/health', (req, res) => {
  res.status(200).send("OK");
});

app.listen(port, () => {
  console.log(`API server listening on port ${port}`);
});
