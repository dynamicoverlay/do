import gql from "graphql-tag";
import { useQuery } from "@apollo/react-hooks";

const userQueryGQL = gql`
    query getUsers {
        user {
            username
            role
            email
        }
    }

`;

export const useUserQuery = () => useQuery(userQueryGQL);
