# System Design Documentation

This folder contains the system design documentation, structured by stage.

## Documentation Stages

- [Requirements](docs/1-requirements.md) - Functional and non-functional requirements.
- [Capacity Estimation](docs/2-capacity-estimation.md) - Back-of-the-envelope calculations for scale.
- [API Design](docs/3-api-design.md) - Endpoint definitions and interfaces.
- [High-Level Design](docs/4-hdl.md) - Core components and interactions.
- [Basic Design Algorithms](docs/5-basic-design-algorithms.md) - Key algorithms and data flows.
- [Database Design](docs/6-database-design.md) - Data models and schema definitions.
- [Cache Design](docs/7-cache-design.md) - Caching strategy and implementation.
- [Detailed Design](docs/8-details-design.md) - Component specifics and optimizations.

First design considered a separate service for generating unique keys for avoiding lock, ACID databases,
collision in hashes and 6 character limit for shortened urls.
You can visit this design [here](https://github.com/alirezaarzehgar/url-shortener/releases/tag/v1.0).
