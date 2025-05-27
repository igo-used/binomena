# Security Guidelines for Binomena Blockchain

## Overview
This document outlines the security measures implemented in the Binomena blockchain and provides guidelines for secure deployment and operation.

## Security Features Implemented

### 1. Authentication & Authorization
- **Admin Key Validation**: All admin endpoints require a valid admin key
- **Environment Variable Support**: Admin key can be set via `ADMIN_KEY` environment variable
- **Audit Logging**: All unauthorized access attempts are logged

### 2. Input Validation
- **Address Format Validation**: All addresses must be 66 characters and start with "AdNe"
- **Amount Validation**: Transaction amounts must be positive and within reasonable limits
- **Percentage Validation**: Token distribution percentages must be valid and sum to 100%
- **Self-Transfer Prevention**: Users cannot transfer tokens to themselves

### 3. Rate Limiting
- **General Endpoints**: 100 requests per minute
- **Transaction Endpoints**: 10 transactions per minute
- **Admin Endpoints**: 5 requests per hour
- **Faucet Endpoints**: 3 requests per hour

### 4. Security Headers
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Referrer-Policy: strict-origin-when-cross-origin`

### 5. CORS Protection
- Configurable CORS middleware to control cross-origin requests

## Security Best Practices

### Environment Variables
Set the following environment variables in production:

```bash
# Required: Admin key for administrative operations
export ADMIN_KEY="your-secure-admin-key-here"

# Optional: Database connection (if using database backend)
export DATABASE_URL="your-database-connection-string"

# Optional: API port (defaults to 8080)
export PORT="8080"
```

### Admin Key Security
- **Never use the default admin key in production**
- Use a strong, randomly generated key (minimum 32 characters)
- Rotate admin keys regularly
- Store admin keys securely (use environment variables or secret management)

### Network Security
- Deploy behind a reverse proxy (nginx, Apache)
- Use HTTPS/TLS in production
- Implement additional firewall rules
- Monitor for suspicious activity

### Database Security
- Use strong database passwords
- Enable database encryption at rest
- Regularly backup database with encryption
- Limit database access to necessary services only

## Deployment Security Checklist

### Pre-Deployment
- [ ] Set strong admin key via environment variable
- [ ] Configure database with strong credentials
- [ ] Review and test all security configurations
- [ ] Ensure no sensitive data in version control

### Production Deployment
- [ ] Deploy behind HTTPS/TLS
- [ ] Configure proper firewall rules
- [ ] Set up monitoring and alerting
- [ ] Enable audit logging
- [ ] Test rate limiting functionality

### Post-Deployment
- [ ] Monitor audit logs regularly
- [ ] Set up automated security scanning
- [ ] Implement backup and recovery procedures
- [ ] Document incident response procedures

## Security Monitoring

### Audit Events to Monitor
- `UnauthorizedAdminAccess`: Invalid admin key attempts
- `TransactionSubmitted`: All transaction activity
- `WalletCreated`: New wallet creation
- `FaucetRequest`: Faucet usage

### Rate Limiting Alerts
Monitor for:
- Repeated rate limit violations from same IP
- Unusual spikes in transaction volume
- Failed authentication attempts

## Incident Response

### Security Incident Types
1. **Unauthorized Admin Access**: Invalid admin key attempts
2. **Rate Limit Abuse**: Excessive requests from single source
3. **Invalid Transaction Patterns**: Suspicious transaction behavior
4. **System Compromise**: Unauthorized system access

### Response Procedures
1. **Immediate**: Block suspicious IP addresses
2. **Short-term**: Rotate admin keys if compromised
3. **Long-term**: Review and update security measures

## Security Updates

### Regular Maintenance
- Update dependencies regularly
- Review and rotate credentials
- Monitor security advisories
- Test backup and recovery procedures

### Security Patches
- Apply security patches promptly
- Test patches in staging environment first
- Document all security-related changes

## Contact Information

For security issues or questions:
- Create a private issue in the repository
- Contact the development team directly
- Follow responsible disclosure practices

## Compliance Notes

This blockchain implementation includes:
- Audit trail capabilities
- Data validation and sanitization
- Access control mechanisms
- Rate limiting and DoS protection

Ensure compliance with relevant regulations in your jurisdiction. 