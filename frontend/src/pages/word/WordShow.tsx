import React, { useEffect, useState } from "react";
import axiosInstance from '../../axiosConfig';
import { useParams, useNavigate, useLocation } from "react-router-dom";


interface Word {
  name: string,
  isRegistered: boolean,
  testCount: number,
  checkCount: number,
  registrationActive: boolean,
  memo: string,
  wordInfos: WordInfo[],
}


  interface WordInfo {
    id: number,
    partOfSpeech: PartOfSpeech
    japaneseMeans: JapaneseMean[]
  }
  
  interface PartOfSpeech {
    id: number,
    name: string,
  }
  interface JapaneseMean {
    id: number,
    name: string,
  }

const WordShow: React.FC = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const [word, setWord] = useState<Word | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  const previousPage = location.state?.page || 1;
  
  const fetchWord = async () => {
    try {
      const response = await axiosInstance.get(`/words/${id}`)
          setWord(response.data)
    } catch (error) {
      console.error("Error fetching word details:", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchWord();
  }, [id]);

  if (loading) {
    return <p>Loading...</p>;
  }

  if (!word) {
    return <p>No word details found.</p>;
  }

  return (
    <div>
      <h1>{word.name}</h1>
      {word.wordInfos.map((info: WordInfo) => (
        <div key={info.id}>
          <p>日本語訳: {info.japaneseMeans.map((japaneseMean: JapaneseMean) => (
            japaneseMean.name
          )).join(', ')}</p>
          <p>品詞: {info.partOfSpeech.name}</p>
        </div>
      ))}
      <p>登録済み: {word.isRegistered ? "はい" : "いいえ"}</p>
      <p>テスト回数: {word.testCount}</p>
      <p>チェック回数: {word.checkCount}</p>
      <p>登録活性: {word.registrationActive ? "はい" : "いいえ"}</p>
      <p>メモ: {word.memo}</p>
      <button onClick={() => navigate('/allwordlist', { state: { page: previousPage } })}>
        一覧に戻る
      </button>
    </div>
  );
};

export default WordShow;
