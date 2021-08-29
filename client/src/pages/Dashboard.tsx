
import { PageHeader, Tag, Button, Statistic, Descriptions, Row, Card } from 'antd';
import { useEffect } from 'react';
import { Route, Switch, useHistory, useRouteMatch } from 'react-router';
import Project from "./Project";
import { useAppDispatch, useAppSelector } from '../app/hooks';

function UserDashboard() {
    // The `state` arg is correctly typed as `RootState` already
    const user = useAppSelector((state: any) => state.user.value)
    console.log("Current User is: ", user);
    const dispatch = useAppDispatch();
    useEffect(() => {
        const loggedInUser = localStorage.getItem("user");
        console.log("logged in:", loggedInUser);
        if (loggedInUser) {
          const foundUser = JSON.parse(loggedInUser);
          dispatch({
            type: "user/set",
            payload: foundUser.Email
          })
        }
    }, []);
    const history = useHistory();
    const redirectToProject = () => {
        history.push("/project/1");
    }
    let match = useRouteMatch();

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
              <Descriptions.Item label="Email">{
                  user ?
                  user
                  : "Lili Qu" 
              }</Descriptions.Item>
              <Descriptions.Item label="Association">
                <a>421421</a>
              </Descriptions.Item>
              <Descriptions.Item label="Creation Time">2017-01-10</Descriptions.Item>
              {/* <Descriptions.Item label="Effective Time">2017-10-10</Descriptions.Item>
              <Descriptions.Item label="Remarks">
                Gonghu Road, Xihu District, Hangzhou, Zhejiang, China
              </Descriptions.Item> */}
            </Descriptions>
            <Row>
              <Statistic title="Status" value="Pending" />
              <Statistic
                title="Balance"
                prefix="$"
                value={568.08}
                style={{
                  margin: '0 45px',
                }}
              />
            </Row>
            <Card
            onClick={() => redirectToProject()}
            style={{ marginTop: 16 }}
            type="inner"
            title="My Project"
            extra={<a href="#">More</a>}
            >
                My special project
            </Card>
          </PageHeader>
        </>
    );
}

export default UserDashboard;
