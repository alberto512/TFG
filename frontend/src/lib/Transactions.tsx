import axios from 'axios';

export const transactions = async (token: string, id: string): Promise<{ code: number; data: any }> => {
  return await axios
    .post(
      process.env.BACKEND_URL || '',
      {
        query: `query AccountById($id: ID!) {
          accountById(id: $id) {
            id,
            iban,
            amount,
            currency,
            transactions {
              id,
              description,
              amount,
              date,
              category {
                id,
                name,
              },
            },
          }
        }`,
        variables: {
          id,
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
      return { code: error.code, data: null };
    });
};

export const updateTransaction = async (
  token: string,
  id: string,
  categoryId: string
): Promise<{ code: number; data: any }> => {
  return await axios
    .post(
      process.env.BACKEND_URL || '',
      {
        query: `mutation UpdateTransaction($id: ID!, $category: ID!) {
          updateTransaction(id: $id, category: $category)
        }`,
        variables: {
          id,
          category: categoryId,
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
