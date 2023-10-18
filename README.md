# generate-container-mod

Dagger module to generate containers using OpenAI

## How to use

Run first query from `examples.gql`: generate container + bash script from prompt.

```shell
export TOKEN="<MY OPENAI TOKEN>"
dagger query script --var token=$TOKEN < examples.gql
```
