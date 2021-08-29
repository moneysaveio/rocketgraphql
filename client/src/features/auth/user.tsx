
import { PageHeader, Tag, Button, Statistic, Descriptions, Row } from 'antd';
import { useAppSelector, useAppDispatch } from '../../app/hooks';

function UserDetails() {
    // The `state` arg is correctly typed as `RootState` already
    const user = useAppSelector((state: any) => state.user.value)
    console.log("Current User is: ", user);
    // const dispatch = useAppDispatch();
    return(
        <>
          <PageHeader
            onBack={() => window.history.back()}
            title="Title"
            tags={<Tag color="blue">Running</Tag>}
            subTitle="This is a subtitle"
            extra={[
              <Button key="3">Operation</Button>,
              <Button key="2">Operation</Button>,
              <Button key="1" type="primary">
                Primary
              </Button>,
            ]}
          >
            <Descriptions size="small" column={3}>
              <Descriptions.Item label="Created">{
                  user ?
                  user
                  : "Lili Qu" 
              }</Descriptions.Item>
              <Descriptions.Item label="Association">
                <a>421421</a>
              </Descriptions.Item>
              <Descriptions.Item label="Creation Time">2017-01-10</Descriptions.Item>
              <Descriptions.Item label="Effective Time">2017-10-10</Descriptions.Item>
              <Descriptions.Item label="Remarks">
                Gonghu Road, Xihu District, Hangzhou, Zhejiang, China
              </Descriptions.Item>
            </Descriptions>
            <Row>
              <Statistic title="Status" value="Pending" />
              <Statistic
                title="Price"
                prefix="$"
                value={568.08}
                style={{
                  margin: '0 32px',
                }}
              />
              <Statistic title="Balance" prefix="$" value={3345.08} />
            </Row>
          </PageHeader>
        </>
    );
}

export default UserDetails;
