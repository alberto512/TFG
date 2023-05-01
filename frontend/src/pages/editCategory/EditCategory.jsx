import React, { useCallback, useEffect, useState } from 'react';
import { useParams } from 'react-router';
import axios from 'axios';
import { useAuth } from 'components/authProvider/AuthProvider';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import './EditCategory.css';

const EditCategory = () => {
  const { id } = useParams();
  const { token } = useAuth();
  const backendUrl = process.env.REACT_APP_BACKEND_URL;
  const [transaction, setTransaction] = useState({});
  const [categories, setCategories] = useState([]);

  const getTransaction = useCallback(() => {
    const getData = async () => {
      const response = await axios.post(
        backendUrl,
        {
          query: `query TransactionById($id: ID!) {
          transactionById(id: $id) {
            id,
            description,
            category {
                id,
                name,
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
      );

      setTransaction(response.data.data.transactionById);
    };

    getData();
  }, [token, backendUrl, id]);

  const getCategories = useCallback(() => {
    const getData = async () => {
      const response = await axios.post(
        backendUrl,
        {
          query: `query { categories() {
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
      );

      setCategories(response.data.data.categories);
    };

    getData();
  }, [token, backendUrl]);

  useEffect(() => {
    getTransaction();
    getCategories();
  }, [getTransaction, getCategories]);

  return (
    <div className='wrapper'>
      {Object.keys(transaction).length === 0 ? (
        <FontAwesomeIcon className='spinner' icon='fa-solid fa-spinner' spin />
      ) : (
        <>
          <div className='title-transaction-wrapper'>
            <span className='title-transaction'>{transaction.description}</span>
            <span className='category-transaction'>
              {transaction.category.name}
            </span>
          </div>
          <div className='scroller'>
            {categories.map((category) => (
              <div key={category.id} className='category-wrapper'>
                <span className='name-category'>{category.name}</span>
              </div>
            ))}
          </div>
        </>
      )}
    </div>
  );
};

export default EditCategory;
