# 🔐 Binomena Environment Configuration Guide

## Security Notice
**NEVER commit real credentials to version control!** All sensitive values should be stored as environment variables or using secure secret management.

## Required Environment Variables

### Database Configuration
```bash
DATABASE_URL=postgresql://username:password@host:port/database_name
DB_PASSWORD=your_secure_database_password
```

### Server Configuration
```bash
PORT=8080                    # Server port (default: 8080)
NODE_ID=your_node_name      # Unique node identifier
DEBUG=false                 # Debug mode (true/false)
GIN_MODE=release            # Gin framework mode (debug/release)
```

### Blockchain Configuration
```bash
GENESIS_BLOCK_HASH=your_genesis_hash
NETWORK_ID=binomena_mainnet
```

## Render.com Deployment

For Render.com deployments, the `render.yaml` file now uses secure database references:

```yaml
envVars:
  - key: DATABASE_URL
    fromDatabase:
      name: binomena-db
      property: connectionString
  - key: DB_PASSWORD
    fromDatabase:
      name: binomena-db
      property: password
```

## Local Development Setup

1. **Create local environment file** (do NOT commit this):
```bash
cp ENV_CONFIG.md .env
# Edit .env with your local database credentials
```

2. **Add .env to .gitignore** (already included):
```bash
echo ".env" >> .gitignore
```

3. **Use environment variables in code**:
```go
dbURL := os.Getenv("DATABASE_URL")
if dbURL == "" {
    log.Fatal("DATABASE_URL environment variable is required")
}
```

## Security Best Practices

✅ **DO:**
- Use environment variables for all secrets
- Use cloud provider secret management (AWS Secrets Manager, etc.)
- Rotate credentials regularly
- Use least-privilege database access
- Enable database SSL/TLS

❌ **DON'T:**
- Commit credentials to Git
- Use default/weak passwords
- Share credentials in chat/email
- Store secrets in code comments
- Use production credentials in development

## Database Security

The database should have:
- Strong password (20+ characters, mixed case, numbers, symbols)
- SSL/TLS encryption enabled
- IP whitelist restrictions
- Regular backups
- Monitoring for suspicious activity

## Recovery

If credentials were exposed:
1. **Immediately** rotate all exposed credentials
2. Check database logs for unauthorized access
3. Review all recent commits for other exposed secrets
4. Update deployment configurations
5. Notify team members if needed 