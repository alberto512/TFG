'use client';
import { useEffect, useState } from 'react';
import Link from 'next/link';
import styles from './Header.module.css';
import { usePathname, useRouter } from 'next/navigation';

const Header = () => {
  const pathname = usePathname();
  let [jwt, setJwt] = useState<string | null>(null);
  const router = useRouter();

  const signOut = () => {
    localStorage.removeItem('jwt');
    if (pathname === '/') {
      router.replace('/login');
    } else {
      router.push('/');
    }
  };

  useEffect(() => {
    setJwt(localStorage.getItem('jwt'));
  }, []);

  return (
    <div className={styles.header}>
      <span className={styles.title}>Bank website</span>
      <div className={styles.navigator}>
        <Link href='/'>
          <span className={`${styles.link} ${pathname == '/' ? styles.selected : ''}`}>Home</span>
        </Link>
        {jwt ? (
          <>
            <Link href='/accounts'>
              <span className={`${styles.link} ${pathname.includes('/accounts') ? styles.selected : ''}`}>Accounts</span>
            </Link>
            <Link href='/stats'>
              <span className={`${styles.link} ${pathname.includes('/stats') ? styles.selected : ''}`}>Stats</span>
            </Link>
            <Link href='/categories'>
              <span className={`${styles.link} ${pathname.includes('/categories') ? styles.selected : ''}`}>Categories</span>
            </Link>
            <Link href='/banks'>
              <span className={`${styles.link} ${pathname.includes('/banks') ? styles.selected : ''}`}>Banks</span>
            </Link>
            <span className={styles.link} onClick={signOut}>
              Sign out
            </span>
          </>
        ) : (
          <>
            <Link href='/login'>
              <span className={`${styles.link} ${pathname.includes('/login') ? styles.selected : ''}`}>Login</span>
            </Link>
            <Link href='/signUp'>
              <span className={`${styles.link} ${pathname.includes('/signUp') ? styles.selected : ''}`}>Sign Up</span>
            </Link>
          </>
        )}
      </div>
    </div>
  );
};

export default Header;
