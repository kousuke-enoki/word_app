import { QuizSettingsType } from '@/types/quiz';
import { QuizStatus } from '@/lib/QuizStatus';

type Props = { setting: QuizSettingsType };

const ResultSettingCard = ({ setting }: Props) => (
  <div className="rs-card">
    <p>登録単語: {QuizStatus.registered[setting.isRegisteredWords as 0 | 1]}</p>
    <p>慣用句　: {QuizStatus.idioms[setting.isIdioms as 0 | 1]}</p>
    <p>特殊単語: {QuizStatus.special[setting.isSpecialCharacters as 0 | 1]}</p>
  </div>
);

export default ResultSettingCard;
