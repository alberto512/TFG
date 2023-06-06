import axios from 'axios';

export const tokenCode = async (token: string, code: string): Promise<{ code: number; data: any }> => {
  return await axios
    .post(
      process.env.BACKEND_URL || '',
      {
        query: `query TokenWithCode($code: String!) {
          tokenWithCode(code: $code)
        }`,
        variables: {
          code: code,
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
