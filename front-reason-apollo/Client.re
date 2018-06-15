/* Create an InMemoryCache */
let inMemoryCache = ApolloInMemoryCache.createInMemoryCache(());

/* Create an HTTP Link */
let httpLink =
 ApolloLinks.createHttpLink(~uri="http://localhost:8080/query", ());

let instance = ReasonApollo.createApolloClient(
 ~link=httpLink,
 ~cache=inMemoryCache,
 ()
);
