import axios from 'axios';

export const balances = async (
  token: string,
  selectedAccounts: string[],
  selectedCategories: string[],
  initDate: number,
  endDate: number
): Promise<{ code: number; data: any }> => {
  return await axios
    .post(
      process.env.BACKEND_URL || '',
      {
        query: `query Balances($accountIds: [ID!]!, $categoryIds: [ID!]!, $initDate: Int!, $endDate: Int!) { 
          balances(accountIds: $accountIds, categoryIds: $categoryIds, initDate: $initDate, endDate: $endDate) {
            amount,
            category {
              name,
            },
          }
        }`,
        variables: {
          accountIds: selectedAccounts,
          categoryIds: selectedCategories,
          initDate: initDate,
          endDate: endDate,
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
