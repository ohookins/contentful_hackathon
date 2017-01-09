## Learnings
* Even when you know the path through the docs in order to get the Management
  API key, there is only one login you can give (email/password) to generate
  a key, and it is unclear what access is given to it. I assume it can be used
  to manage in any space you have access to, but don't know for sure.
* Contentful Go SDK is still very rough around the edges :)
* Even if it weren't, the content structure is complex enough to make it a
  challenge to generate something very simple with existing content. Too many
  additional metadata fields compared to how much content you have.
* Go SDK currently lacks any CMA methods :(
* CMA docs don't make it easy to understand how to construct a request in code.
  It is _kinda_ human understandable but the way it is written doesn't map well
  to written code.
* You can't easily find the ID of the Content Type in the Web UI.
* "fields" structure in CMA entries is case sensitive :(

## Running
You need to set these three environment variables:
```
CFHACK_CMA_TOKEN
CFHACK_CDA_TOKEN
CFHACK_SPACE_ID
```
