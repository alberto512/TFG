import axios from 'axios';

export const create = async (token: string, name: string): Promise<{ code: number; data: any }> => {
  return await axios
    .post(
      process.env.BACKEND_URL || '',
      {
        query: `mutation CreateCategory($name: String!) {
          createCategory(name: $name) {
            id,
            name,
          }
        }`,
        variables: {
          name,
        },
      },
      {
        headers: {
          Authorization: token,
          withCredentials: true,
        },
      }
    )
    .then((response) => {
      return { code: 200, data: response.data };
    })
    .catch((error) => {
      console.log(error);
      try {
        if (error.response.data === 'Invalid token\n') {
          return { code: 401, data: null };
        }
        return { code: error.response.status, data: null };
      } catch (error) {
        return { code: 500, data: null };
      }
    });
};

export const categories = async (token: string): Promise<{ code: number; data: any }> => {
  return await axios
    .post(
      process.env.BACKEND_URL || '',
      {
        query: `query { categories() 
          {
            id,
            name,
          }
        }`,
      },
      {
        headers: {
          Authorization: token,
          withCredentials: true,
        },
      }
    )
    .then((response) => {
      return { code: 200, data: response.data };
    })
    .catch((error) => {
      console.log(error);
      try {
        if (error.response.data === 'Invalid token\n') {
          return { code: 401, data: null };
        }
        return { code: error.response.status, data: null };
      } catch (error) {
        return { code: 500, data: null };
      }
    });
};
