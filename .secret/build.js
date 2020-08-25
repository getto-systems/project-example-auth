const fs = require("fs");
const path = require("path");

function adminSecret(cwd) {
  try {
    return JSON.parse(fs.readFileSync(path.join(cwd, "admin.json")));
  } catch (err) {
    console.error(`failed to create admin secret: ${err}`);
  }
}

function cookieSecret(cwd) {
  try {
    return JSON.parse(fs.readFileSync(path.join(cwd, "cookie.json")));
  } catch (err) {
    console.error(`failed to create cookie secret: ${err}`);
  }
}

function ticketSecret(cwd) {
  try {
    return {
      private_key: fs.readFileSync(path.join(cwd, "ticket/ecdsa-p521-private.pem")).toString(),
      public_key: fs.readFileSync(path.join(cwd, "ticket/ecdsa-p521-public.pem")).toString(),
    }
  } catch (err) {
    console.error(`failed to create ticket secret: ${err}`);
  }
}

function apiSecret(cwd) {
  try {
    return {
      private_key: fs.readFileSync(path.join(cwd, "api/ecdsa-p521-private.pem")).toString(),
    }
  } catch (err) {
    console.error(`failed to create api secret: ${err}`);
  }
}

function cloudfrontSecret(cwd) {
  try {
    const secret = JSON.parse(fs.readFileSync(path.join(cwd, "cloudfront.json")));
    secret.private_key = fs.readFileSync(path.join(cwd, "cloudfront/pk.pem")).toString();
    return secret;
  } catch (err) {
    console.error(`failed to create cloudfront secret: ${err}`);
  }
}

function buildSecret(cwd) {
  fs.writeFileSync(path.join(cwd, "secret.json"), JSON.stringify({
    admin: adminSecret(cwd),
    cookie: cookieSecret(cwd),
    ticket: ticketSecret(cwd),
    api: apiSecret(cwd),
    cloudfront: cloudfrontSecret(cwd),
  }));
}

buildSecret(process.cwd());
