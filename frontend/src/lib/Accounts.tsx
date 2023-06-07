import axios from 'axios';

export const refresh = async (token: string): Promise<{ code: number; data: any }> => {
  console.log(token);
  return await axios
    .post(
      process.env.BACKEND_URL || '',
      {
        query: `mutation { refreshBankData }`,
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

export const accounts = async (token: string): Promise<{ code: number; data: any }> => {
  console.log(token);
  return await axios
    .post(
      process.env.BACKEND_URL || '',
      {
        query: `query { 
          accounts() {
            id,
            iban,
            name,
            currency,
            amount,
            bank,
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

export const fullAccounts = async (token: string): Promise<{ code: number; data: any }> => {
  return await axios
    .post(
      process.env.BACKEND_URL || '',
      {
        query: `query { 
          accounts() {
            id,
            iban,
            transactions {
              id,
              date,
              amount,
              category {
                id,
                name,
              },
            },
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
