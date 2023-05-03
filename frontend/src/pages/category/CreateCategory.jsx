import React, { useState } from 'react';
import { useAuth } from 'components/authProvider/AuthProvider';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import './CreateCategory.css';

const CreateCategory = () => {
  const navigate = useNavigate();
  const { token } = useAuth();
  const [name, setName] = useState('');
  const [loading, setLoading] = useState(false);
  const backendUrl = process.env.REACT_APP_BACKEND_URL;

  const handleSubmit = async () => {
    if (name === '') {
      return;
    }

    setLoading(true);

    await axios.post(
      backendUrl,
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
    );

    setLoading(false);

    navigate(-1);
  };

  return (
    <div className='wrapper'>
      <div className='form-container'>
        <div className='form-wrapper'>
          <span className='form-label'>Name</span>
          <input
            className='input-label'
            type='text'
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </div>
      </div>
      {loading ? (
        <FontAwesomeIcon className='spinner' icon='fa-solid fa-spinner' spin />
      ) : (
        <div className='btn-submit' onClick={handleSubmit}>
          <span>Create category</span>
        </div>
      )}
    </div>
  );
};

export default CreateCategory;
