'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import axios from 'axios';
import Button from '@elements/Button/Button';
import styles from './CategoriesPage.module.css';
import InputText from '@elements/InputText/InputText';

const CategoriesPage = () => {
  const router = useRouter();
  const [categoryName, setCategoryName] = useState('');
  const [categories, setCategories] = useState<Category[]>([]);
  const [createModal, setCreateModal] = useState(false);

  const disableModal = () => {
    setCreateModal(false);
    setCategoryName('');
  };

  const createCategory = () => {
    axios
      .post(
        '/api/categories',
        { name: categoryName },
        {
          headers: {
            Authorization: localStorage.getItem('jwt'),
          },
        }
      )
      .then((_response) => {
        axios
          .get('/api/categories', {
            headers: {
              Authorization: localStorage.getItem('jwt'),
            },
          })
          .then((response) => {
            setCategories(response.data);
            disableModal();
          })
          .catch((error) => {
            console.log(error);
            if (error.response.status === 401) {
              localStorage.removeItem('jwt');
              router.push('/login');
            }
          });
      })
      .catch((error) => {
        console.log(error);
        if (error.response.status === 401) {
          localStorage.removeItem('jwt');
          router.push('/login');
        }
      });
  };

  useEffect(() => {
    axios
      .get('/api/categories', {
        headers: {
          Authorization: localStorage.getItem('jwt'),
        },
      })
      .then((response) => {
        setCategories(response.data);
      })
      .catch((error) => {
        console.log(error);
        if (error.response.status === 401) {
          localStorage.removeItem('jwt');
          router.push('/login');
        }
      });
  }, []);

  return (
    <>
      {createModal && (
        <div className={styles.createModal} onClick={disableModal}>
          <div className={styles.modalWrapper} onClick={(event) => event.stopPropagation()}>
            <span className={styles.modalTitle}>Create a category</span>
            <InputText iconPath='/svg/category.svg' onChange={setCategoryName} />
            <Button label={'Create'} onClick={createCategory} />
          </div>
        </div>
      )}
      <div className={`${styles.categoriesPage} ${createModal && styles.categoriesPageOpacity}`}>
        <div className={styles.titleWrapper}>
          <div className={styles.mainTitleWrapper}>
            <span className={styles.title}>Categories</span>
            <div className={styles.buttonWrapper}>
              <Button label={'Create category'} onClick={() => setCreateModal(true)} />
            </div>
          </div>
          <div className={styles.keysWrapper}>
            <span className={styles.key}>Name</span>
          </div>
        </div>
        <div className={styles.categoriesWrapper}>
          {categories.map((category, index, categories) => {
            return (
              <div
                key={category.id}
                className={`${styles.category} ${index + 1 === categories.length ? styles.lastCategory : ''}`}
              >
                <span className={`${styles.item} ${styles.name}`}>{category.name}</span>
              </div>
            );
          })}
        </div>
      </div>
    </>
  );
};

export default CategoriesPage;
