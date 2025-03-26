# Simple Bank - Front End and Back End

## Introduction

This branch is an extension of the original repository, maintained by [TechSchool](https://github.com/techschool) and [Quang Pham](https://github.com/phamlequang)

The purpose of this branch is to explore and build upon their work [Backend Master Class](https://github.com/techschool/simplebank) and [Front End Crash](https://www.youtube.com/playlist?list=PLy_6D98if3UI3rsFRTHM1LMtVprYMp-GT) Courses by introducing new ideas, changes, features, improvements, or exploring new directions for this project.

Some tech I have been learning through this:
- sqlcgen
- Vue.js
- protobuffer
- gRPC
- gin
- paseto
- migration
- swagger

## Next Steps

- Future Ideas:
    - Adding terraform files to deploy the AWS resources to apply and destroy simplebank app deployment : RDS, EKS, Route53 (ingress) to follow on from the course and suggest any improvements.  It costs money to leave these resources running, so it would be cool to shutdown after demo, development and test.
        - I would like to deploy the necessary resources when I am using them, terraform destroy when I am done without using the console.
        - Also add a lower cost deployment for development and test purposes
    - Adding component tests: 
        - extending the current docker compose file to include a container running a number of component tests
        - maybe include some basic UI tests with Selenium
    - Having completed the Front End Crash Course:
        - Explore and extend the VUE js code to include other workflows:
            - Create User
            - Balance Transfers
        - Add some tests
        - Add front end AWS deployment to the terraform

- Documentation Updates:
    - I aim to share any new instructions here

### So far

- Branches:
    - **Main**: Any diversions or additions to the original fork and Udemy which intended to be shared with the author for future considerations.
        - Terraform : Adding initial terraform files for: secrets manager, RDS (will need to change subnet setup) and ECR
        - This code is not yet exhaustively complete with unit tests and app code, so add more features and tests as I explore the frontend and aws resources.
    - **Udemy**: Following the Udemy course step by step and adding the code commits after each section I complete.  
        - This code is not yet exhaustively complete with unit tests and app code.
        - Code from following the youtube Frontend Crash Course included.
        - Completed 99% the Udemy Course.
            - Remaining : Automatic deploy to Kubernetes with Github Action
                - I dont intend to use the github actions frequently and I would like to get the deployment automated with terraform before I complete this one.
            - App deployed after using AWS Console to create EKS and node group, RDS, Elastic Cache, Security Groups and both the grpc and http calls work for loginUser and createUser from postman.
    - **Master**: Fork of the Backend Master Class code
