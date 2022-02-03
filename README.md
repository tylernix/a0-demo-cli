# Personal CLI Tool To Help With Auth0 Demos

I wanted a tool to automate some of the repetitive tasks I was talking about while performing demos of the Auth0 product. I also wanted to learn Golang in the process.


## Demo

![a0-demo-user-import](https://user-images.githubusercontent.com/67964959/152444909-aa299d50-402d-443d-b854-011afeb78bd8.gif)

## Environment Variables

To run this project, you will need to add the following environment variables to a `local.env` file in the root of the project.

`AUTH0_DOMAIN`

`AUTH0_CLIENT_ID`

`AUTH0_CLIENT_SECRET`

These environment variables can be obtained by creating a M2M application in your Auth0 tenant. Make sure it is given the permissions to call call the Auth0 Management API.
