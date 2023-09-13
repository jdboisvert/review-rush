# Review Rush

A Slack bot used to show who the top reviewers were in GitHub for a given day.

## Setting Up and Running the Project with Slack Integration

### Prerequisites:

1. **Go**: Ensure you have Go installed. [Download and install Go](https://golang.org/dl/).
2. **Slack Token**: Obtain a Slack token (by creating an app) to enable posting messages. [Follow Slack's documentation](https://api.slack.com/tutorials/tracks/getting-a-token) to get one.
3. **Slack Channel**: Create a Slack channel to post messages to. [Follow Slack's documentation](https://slack.com/help/articles/201402297-Create-a-channel) to create one.
4. **GitHub Token**: Obtain a GitHub token to enable fetching data from GitHub. [Follow GitHub's documentation](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token) to get one. Ensure it has the following scopes `public_repo, read:project, read:user, repo:status, user:email`. When making this project I simply used a Classic Token. Ideally just using a Fine-grained personal access token would be better since scoping by organization or user would be more secure.

### Quick Start:

1. **Clone the Repository**:
    ```bash
    git clone git@github.com:jdboisvert/review-rush.git
    cd review-rush
    ```

2. **Set Environment Variables**:
    ```bash
    cp .env.example .env
    ```
    This project makes use of a .env file you need to copy .env.example to .env and fill in the values.

3. **Build the Project for Quick Use using /cmd**:
    ```bash
    cd ./cmd
    go build -o ../review-rush.exe
    ```

4. **Run the Built Executable**:
    ```bash
    cd ..
    ./review-rush.exe
    ```

    You can skip the build process and run the project directly using `go run ./cmd/main.go` from the root of the project as well. 

5. The application should now be running and posting to Slack as intended!
   ![image](https://github.com/jdboisvert/review-rush/assets/40838156/214f2ac2-f9a4-4c51-954c-2e3dfaaef466)


### Development Configuration:

If you need more configuration options or to understand deeper aspects of the integration, [visit this detailed guide](./docs/development.md).



