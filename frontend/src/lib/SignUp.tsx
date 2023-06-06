import axios from 'axios';

export const signUp = async (username: string, password: string): Promise<{ code: number; data: any }> => {
  return await axios
    .post(process.env.BACKEND_URL || '', {
      query: `mutation CreateUser($username: String!, $password: String!, $role: Role!) {
        createUser(username: $username, password: $password, role: $role) {
          id,
          username,
          password,
        }
      }`,
      variables: {
        username,
        password,
        role: 'USER',
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
