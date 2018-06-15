/* Create a GraphQL Query by using the graphql_ppx */
module GetPokemon = [%graphql {|
    query getPokemon($name: String!){
        pokemon(name: $name) {
            name
        }
    }
  |}];

  module GetPokemonQuery = ReasonApollo.CreateQuery(GetPokemon);

  let make = (_children) => {
  /* ... */,
  render: (_) => {
    let pokemonQuery = GetPokemon.make(~name="Pikachu", ());
    <GetPokemonQuery variables=pokemonQuery##variables>
      ...(({result}) => {
        switch result {
           | Loading => <div> (ReasonReact.string("Loading")) </div>
           | Error(error) => <div> (ReasonReact.string(error)) </div>
           | Data(response) => <div> (ReasonReact.string(response##pokemon##name)) </div>
        }
      })
    </GetPokemonQuery>
  }
  }
