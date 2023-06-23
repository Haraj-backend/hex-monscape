# Project Methodology

This document contains guidelines on how to implement [Hexagonal Architecture](../reference/hex-architecture.md) in Solutions Team projects.

1. Understand business requirements. This is very important because it gives us context on what we are building.
2. Validate your understanding by creating document that contains high level overview of the application. Usually we use `README.md` for this.
3. Write down the expected use cases for the project. It will be much better if we can also provide use case diagrams when writing this.
4. Write API specification for the project based on the use cases we wrote in `step 3`.
5. Identify the usage context of your API, this is to spot out your `Core` componentes for your project.
6. Create class diagram for the core components. By creating it, it will be easier for you to structure your code. We can use [UMLet](https://www.umlet.com/) to create the class diagram. It is okay to skip this test if the `Core` components are not that complex.
7. Start writing your code & make any necessary adjustment to the docs.
8. Create necessary tests for your project.
9. Dockerize your project so it could be run easily by other team members.
10. Open pull request & ask review from your teammates.