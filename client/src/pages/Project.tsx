
import { PageHeader, Tag, Button, Statistic, Descriptions, Row, Card, Tabs } from 'antd';
import { useEffect } from 'react';
import { Route, Switch, useHistory, useRouteMatch } from 'react-router';
import Project from "./Project";
import { useAppDispatch, useAppSelector } from '../app/hooks';
import { useParams } from 'react-router-dom';
import axios from 'axios';
import { AppleOutlined, AndroidOutlined } from '@ant-design/icons';

const { TabPane } = Tabs;

function UserProject() {

    const getProjectDetails = (id: string) => {
        const API_URL = `${process.env.REACT_APP_API_ENDPOINT}/project/${id}`;
        axios.get(API_URL)
            .then(response => {
                console.log(response);
            })
    }
    let { projectId } = useParams();
    const details = getProjectDetails(projectId);
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
            <Tabs defaultActiveKey="2">
                <TabPane
                tab={
                    <span>
                    <AppleOutlined />
                    PostgreSQL
                    </span>
                }
                key="1"
                >
                Tab 1
                </TabPane>
                <TabPane
                tab={
                    <span>
                    <AndroidOutlined />
                    Hasura
                    </span>
                }
                key="2"
                >
                Tab 2
                </TabPane>
            </Tabs>
          </PageHeader>
        </>
    );
}

export default UserProject;
