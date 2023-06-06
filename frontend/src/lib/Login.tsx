import axios from 'axios';

export const login = async (username: string, password: string): Promise<{ code: number; data: any }> => {
  return await axios
    .post(process.env.BACKEND_URL || '', {
      query: `mutation Login($username: String!, $password: String!) {
        login(username: $username, password: $password)
      }`,
      variables: {
        username,
        password,
      },
    })
    .then((response) => {
      return { code: 200, data: response.data };
    })
    .catch((error) => {
      console.log(error);
      try {
        return { code: error.response.status, data: null };
      } catch (error) {
        return { code: 500, data: null };
      }
    });
};
