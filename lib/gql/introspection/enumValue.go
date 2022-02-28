package introspection
type EnumValue struct{
	Name 				string 	`gql:"name=name"`
	Description 		*string	`gql:"name=description"`
	IsDeprecated		bool	`gql:"name=isDeprecated"`
	DeprecationReason 	*string	`gql:"name=deprecationReason"`
}