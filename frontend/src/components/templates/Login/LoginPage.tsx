'use client';
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Image from 'next/image';
import axios from 'axios';
import InputText from '@elements/InputText/InputText';
import Button from '@elements/Button/Button';
import styles from './LoginPage.module.css';

const LoginPage = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const router = useRouter();

  const handleLogin = () => {
    axios
      .post('/api/login', {
        username,
        password,
      })
      .then((response) => {
        localStorage.setItem('jwt', response.data);
        router.push('/accounts');
      })
      .catch((error) => {
        console.log(error);
      });
  };

  return (
    <div className={styles.loginPage}>
      <div className={styles.imageWrapper}>
        <Image src={'/svg/login.svg'} alt='Login image' height={500} width={500} />
      </div>
      <div className={styles.loginWrapper}>
        <div className={styles.item}>
          <InputText iconPath='svg/user.svg' onChange={setUsername} />
        </div>
        <div className={styles.item}>
          <InputText iconPath='svg/password.svg' onChange={setPassword} type='password' />
        </div>
        <div className={styles.button}>
          <Button label='Login' onClick={handleLogin} />
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
