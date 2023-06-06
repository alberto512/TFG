'use client';
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Image from 'next/image';
import axios from 'axios';
import InputText from '@elements/InputText/InputText';
import Button from '@elements/Button/Button';
import styles from './SignUpPage.module.css';

const SignUpPage = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const router = useRouter();

  const handleSignUp = () => {
    axios
      .post('/api/signUp', {
        username,
        password,
      })
      .then((_response) => {
        axios
          .post('/api/login', {
            username,
            password,
          })
          .then((response) => {
            localStorage.setItem('jwt', response.data);
            router.push('/banks');
          })
          .catch((error) => {
            console.log(error);
          });
      })
      .catch((error) => {
        console.log(error);
      });
  };

  return (
    <div className={styles.signUpPage}>
      <div className={styles.imageWrapper}>
        <Image src={'/svg/signUp.svg'} alt='Sign Up image' height={500} width={500} />
      </div>
      <div className={styles.signUpWrapper}>
        <div className={styles.item}>
          <InputText iconPath='svg/user.svg' onChange={setUsername} />
        </div>
        <div className={styles.item}>
          <InputText iconPath='svg/password.svg' onChange={setPassword} type='password' />
        </div>
        <div className={styles.button}>
          <Button label='Sign up' onClick={handleSignUp} />
        </div>
      </div>
    </div>
  );
};

export default SignUpPage;
