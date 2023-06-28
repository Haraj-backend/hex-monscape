# Project Methodology

This document contains guidelines on how to implement [Hexagonal Architecture](../reference/hex-architecture.md) in Solutions Team projects.

1. Understand business requirements. This is very important because it gives us context on the solution we are trying to build.
2. Validate our understanding by creating document that contains high level overview of the application. Usually we use `README.md` file for this.
3. Write down the expected use cases for the project. It will be much better if we can also provide use case diagrams when writing them.
4. Write API specification for the project based on the use cases we wrote in `step 3`.
5. Identify the usage context of our API. This is to spot out our [`Core`](./hex-architecture.md#core) components.
6. Create class diagram for the [`Core`](./hex-architecture.md#core) components. By creating it, it will be easier for us to structure our code. We can use [UMLet](https://www.umlet.com/) to create the class diagram. It is okay to skip this step if our [`Core`](./hex-architecture.md#core) components seems not that complex.
7. Start writing our code & make any necessary adjustment to the docs created in previous steps.
8. Create necessary tests for our project. Remember our objective is not to get `100%` coverage, but to cover the most important scenarios in our project.
9. Dockerize our project so it could be run easily by other team members. Also create `Makefile` to simplify the process of running it.
10. Open pull request & ask review from your teammates.