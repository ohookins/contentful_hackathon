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
  * Actually the SDK is harder to use than raw Go and the APIs, currently.
* CMA docs don't make it easy to understand how to construct a request in code.
  It is _kinda_ human understandable but the way it is written doesn't map well
  to written code.
* You can't easily find the ID of the Content Type in the Web UI.
* "fields" structure in CMA entries is case sensitive :(
* To publish an entry you call the "published" endpoint which doesn't make much
  sense.
* You can't create and publish an entry in the same API call :(
* You MUST supply a version ID when publishing, it won't just be the most
  recently created ID.
* Is it possible to set the content type in an entry query via a header?!?
  Inconsistent use of headers and query parameters...
* Can't see the short field names in the Content Type screen of the UI.
* Go is not a great language for dynamic type structures that you might need for
  Contentful apps.

## Running
You need to set these three environment variables:
```
CFHACK_CMA_TOKEN
CFHACK_CDA_TOKEN
CFHACK_SPACE_ID
PORT
```
