import { Card } from '@/components/ui/card'
import { QuizStatus } from '@/lib/QuizStatus'
import { QuizSettingsType } from '@/types/quiz'

type Props = { setting: QuizSettingsType }

const ResultSettingCard = ({ setting }: Props) => (
  <Card className="mb-6 p-4">
    <h2 className="mb-3 text-sm font-semibold opacity-80">出題条件</h2>
    <dl className="grid gap-3 sm:grid-cols-3">
      <div>
        <dt className="text-xs opacity-70">登録単語</dt>
        <dd className="text-sm font-medium">
          {QuizStatus.registered[setting.isRegisteredWords as 0 | 1]}
        </dd>
      </div>
      <div>
        <dt className="text-xs opacity-70">慣用句</dt>
        <dd className="text-sm font-medium">
          {QuizStatus.idioms[setting.isIdioms as 0 | 1]}
        </dd>
      </div>
      <div>
        <dt className="text-xs opacity-70">特殊単語</dt>
        <dd className="text-sm font-medium">
          {QuizStatus.special[setting.isSpecialCharacters as 0 | 1]}
        </dd>
      </div>
    </dl>
  </Card>
)

export default ResultSettingCard
