This documentation clarified functional and non-functional requirements.

The final product is URL shortener. Users can create short links and
that links will redirect to offered links. No need to register and login and
everyone should be able to use service.

# Functional requirements
- Create a short link from offered orbitrary link
- Redirect link after sending request to it (3xx code)
- Links will expire after specified time

# Non-Functional requirements
- The system should be highly available with 99.99% uptime
- URL redirection should happen in less than 50ms
- System should be scalable for higher read/write loads
- Short links should be unique
- DDOS protection and do not allow attackers to reserve all keys or full bandwidth:
  - Rate limiting is needed
  - Block suspicious IPs
- Cleanup expired data to free space for new short links
- Telemetry:
  - Calculate how many times links redirected, success and failure rates
  - Performance metrics
  - Detect sucpicious traffic and IPs
