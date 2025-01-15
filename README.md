# Simple Bank - Front End and Back End

## Introduction
=======

This repository contains the codes of the [Backend Master Class](https://bit.ly/backendmaster) course by [TECH SCHOOL](https://bit.ly/m/techschool).

This branch is an extension of the original repository, maintained by [TechSchool](https://github.com/techschool) and [Quang Pham](https://github.com/phamlequang)

The purpose of this branch is to explore and build upon their work [Backend Master Class](https://github.com/techschool/simplebank) and [Front End Crash](https://www.youtube.com/playlist?list=PLy_6D98if3UI3rsFRTHM1LMtVprYMp-GT) Courses by introducing new ideas, changes, features, improvements, or exploring new directions for this project.

## Next Steps

- Future Ideas:
    - Adding terraform files to deploy the AWS resources to apply and destroy simplebank app deployment : RDS, EKS, Route53 (ingress) to follow on from the course and suggest any improvements.  It costs money to leave these resources running, I would like to deploy the necessary resources when I am using them, destroy when I am done without using the console.
    - Adding component tests, extending the current docker compose file to inc lude tests
    
- Complete the Udemy Course:
    - Complete the course: EKS, manage K8sm, deploy web app, auto TLS

- Having completed the Front End Crash Course:
    - Explore and extend the VUE js code to include other workflows:
        - Create User
        - Balance Transfers
    - Add some tests

- Documentation Updates:
    - I aim to share any new instructions here

### So far

- Branches:
    - Main: Any diversions or additions to the original fork and Udemy which intended to be shared with the author for future considerations.
        - Terraform : Adding initial terraform files for: secrets manager, RDS (will need to change subnet setup) and ECR
        - This code is not yet exhaustively complete with unit tests and app code, so add more features and tests as I explore the frontend and aws resources.
    - Udemy: Following the Udemy course step by step and adding the code after each section I complete.  
        - This code is not yet exhaustively complete with unit tests and app code.
        - Code from following the youtube Frontend Crash Course included.
    - Master: Fork of the Backend Master Class code
