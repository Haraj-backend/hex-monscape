# Project Methodology

This document contains guidelines on implementing [Hexagonal Architecture](../reference/hex-architecture.md) to a new Solutions Team project.

1. Understand the business requirements. This is very important because it gives us context on the solution we are trying to build.
2. Validate our understanding by creating a document containing a high-level application overview. Usually we use `README.md` file for this.
3. Write down the expected use cases for the project. It would be better if we can also create use case diagrams when writing them.
4. Write API specification for the project based on the use cases we wrote in `step 3`.
5. Spot out our `Core` components following the guidelines mentioned in [here](./hex-architecture.md#core).
7. While implementing the code, make any necessary adjustments to the docs created in previous steps. This is to make our docs stay up to date.
8. Create necessary tests for our project. Remember, our objective is not to get `100%` coverage but to cover the most essential scenarios in our project.
9. Dockerize our project so it can be run easily by other team members. Also create `Makefile` to simplify the process of running it.
10. Open a pull request & ask review from your teammates.