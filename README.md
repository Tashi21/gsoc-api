# GSoC-API

A RESTful API that presents all the past organizations that have been a part of GSoC. This is a personal project to learn the basics of REST and to learn Golang. This might be used as a reference for my [GSoC](https://github.com/Tashi21/gsoc) project.

## [template.env](https://github.com/Tashi21/gsoc-api/blob/main/template.env)

This is a template .env file with all the key names being used in the project. Please add your own details as needed.

## Endpoints

### 1. /orgs

A GET request to get all the organizations.

### 2. /orgs/{org_name}

A GET request to get the details of a particular organization.

### 3. /orgs/{org_name}

A PATCH request to edit a particular organization.

### 4. /orgs/{org_name}

A DELETE request to delete a particular organization.
