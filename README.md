# Nuance Retrieval Service 

## Overview 

The Nuance Retrieval Service is exposes a RAG application over a tcp server


// query = "Which training method should I use for sentence transformers when " +
// 	"I only have pairs of related sentences?"

// res = openai.Embedding.create(
// 	input=[query],
// 	engine=embed_model
// )

// # retrieve from Pinecone
// xq = res['data'][0]['embedding']

// # get relevant contexts (including the questions)
// res = index.query(vector=xq, top_k=2, include_metadata=True)

// query = (
// 	"Which training method should I use for sentence transformders when " +
// 	"I only have pairs of related sentences?"
// )

// res = openai.Embedding.create(
// 	input=[query],
// 	engine=embed_model
// )

// # retrieve from Pinecone
// xq = res['data'][0]['embedding']

// # get relevant contexts (including the questions)
// res = index.query(vector=xq, top_k=2, include_metadata=True)

// # first we retrieve relevant items from Pinecone
// query_with_contexts = retrieve(query)
// query_with_contexts

// # then we complete the context-infused query
// complete(query_with_contexts)

// pc.delete_index(index_name)