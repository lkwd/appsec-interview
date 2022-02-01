# Docu-me

This repository contains an early prototype of Docu-me, a document storage and management REST API. Right now it only allows
very basic functionality - authenticated users can view files and an administrator can rename them. We hope to expand
this in the future.

Unfortunately, shortly after deploying a preview instance we noticed some very strange activity. We also received the 
following anonymous message:

> Your devs made some serious mistakes. I can completely take over Docu-me via RCE, without even needing a user account.
> Pay 1 BTC to wallet_address_here and I'll help you fix this.


Unsurpisingly, we don't want to pay this mysterious and possibly malevolent actor. Instead, we hope you'll be able to
help us work out what they might have found and help us fix it. We know you're very busy, but we hope this won't take
more than 2-4 hours of your time.

## Your task

### Step 1

* Determine how an unauthenticated user might be able to gain RCE on the server
* Provide a written report for our developers so they can understand how this is possible (this need only be a few paragraphs, but should sufficiently clear)
* Clone the repository and fix all the issues _directly leading_ to that RCE without breaking the core functionality of:
  * Allowing authenticated users to read documents
  * Allowing the `admin` user to rename documents
  * Allowing anonymous users to login
* Document any additional security issues you find in the process

### Step 2

If you have any time left, we'd greatly appreciate it if you could add a new endpoint that allows the `admin` user to
delete files via the REST API. Add a `deleteHandler` to `main.go` and appropriate logic to `fileops` to allow this.

### Step 3

After completing step 1 (and optionally step 2), please either:

1. Upload it to your own Github account and grant access to `thomas-welch` (Preferred)
2. Zip up the repository (including the .git folder) and email it to `thomas.welch@lockwood-publishing.com`

## Building and running

* this code is known to compile under go 1.17.2 darwin/amd. It should also run under linux or WSL
* compilation should be as easy as running `go build -o docume main.go`. This will create a binary called `docume` in the current directory that can be run using `./docume`
* the application relies on a 32 byte secret that should be written to `secret.txt` in the same directory the application is executed from - you can generate one by running `cat /dev/urandom | env LC_ALL=C tr -dc 'a-zA-Z0-9' | head -c 32 > secret.txt`

## Notes

Only files in the `./files` directory should be accessible to or modifiable by any users (including `admin`).

As this is pre-alpha software we're aware there are some things that need to be done before we're ready for production.
You may spot some of the following:

* Lack of HTTPS/TLS
* Missing security headers
* Hardcoded admin password
* An in-memory stub instead of a real user database
* No way of creating new users

We are already aware of these issues but don't believe they are what the mysterious actor was referring to. You may choose
to address some of these issues in your report/patch, but they are not the main goal of this exercise.

## Examples

There are 3 REST endpoints currently available:

### 1. View 

Allows authenticated users to view the files

* URL: `http://localhost:9555/view`
* Method: `GET`
* Headers:
  * `Authorization: Bearer <token_here>` 

```
 NO REQUEST BODY
```

### 2. Rename

Allows an admin user to rename a file

* URL: `http://localhost:9555/rename`
* Method: `POST`
* Headers:
  * `Authorization: Bearer <token_here>`
  * `Content-Type: application/json` 
  
```json
{
  "old": "oldfile.txt",
  "new": "newfile.txt"
}
```

### 3. Login

Allows anonymous users to log in

* URL: `http://localhost:9555/login`
* Method: `POST`
* Headers:
  * `Authorization: Bearer <token_here>`
  * `Content-Type: application/json` 

```json
{
	"name": "MY_NAME",
	"password": "MY_PASSWORD"
}
```